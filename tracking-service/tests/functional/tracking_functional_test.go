// tests/functional/tracking_functional_test.go
// FUNCTIONAL TEST untuk Tracking & Status Service.
//
// Build tag memastikan test ini HANYA dijalankan secara eksplisit:
//   go test -v -tags=functional -count=1 ./tests/functional/...
//
// Tanpa tag itu, test ini akan diskip otomatis saat unit test biasa.
//
// Test ini MENGAKSES database nyata (MongoDB + Redis) yang dijalankan
// oleh Docker Compose di CI/CD pipeline.

//go:build functional
// +build functional

package functional_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"tracking-service/internal/cache"
	"tracking-service/internal/domain"
	"tracking-service/internal/handler"
	"tracking-service/internal/repository"
	"tracking-service/internal/service"
)

// =========================================================
// TrackingFunctionalSuite — Test Suite utama
// Menggunakan testify/suite untuk setup/teardown yang terstruktur
// =========================================================

type TrackingFunctionalSuite struct {
	suite.Suite
	svc     *service.TrackingService
	handler *handler.TrackingHandler
	ctx     context.Context
	awbSeq  int // counter untuk generate AWB unik per test
}

// SetupSuite dijalankan SEKALI sebelum semua test.
// Koneksi ke MongoDB dan Redis dibuat di sini.
func (s *TrackingFunctionalSuite) SetupSuite() {
	s.ctx = context.Background()

	// --- Koneksi MongoDB ---
	mongoURI := os.Getenv("TEST_MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://testuser:testpass@localhost:27017/tracking_test?authSource=admin"
	}
	mongoDBName := os.Getenv("TEST_MONGO_DB")
	if mongoDBName == "" {
		mongoDBName = "tracking_test"
	}

	db, err := repository.ConnectMongoDB(mongoURI, mongoDBName)
	s.Require().NoError(err, "Gagal koneksi ke MongoDB test")

	// --- Koneksi Redis ---
	redisAddr := os.Getenv("TEST_REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6380" // Port berbeda agar tidak bentrok dengan Redis dev
	}

	redisClient, err := cache.ConnectRedis(redisAddr, "", 1) // DB index 1 untuk test
	s.Require().NoError(err, "Gagal koneksi ke Redis test")

	// Gunakan implementasi NYATA (bukan mock)
	repo := repository.NewMongoTrackingRepository(db)
	redisCache := cache.NewRedisTrackingCache(redisClient)

	// Buat index MongoDB untuk performa
	err = repo.EnsureIndexes(s.ctx)
	s.Require().NoError(err, "Gagal membuat MongoDB index")

	// Kafka diganti NoOp agar functional test tidak butuh Kafka nyata
	noopKafka := &NoOpKafkaProducer{}

	s.svc = service.NewTrackingService(repo, redisCache, noopKafka)
	s.handler = handler.NewTrackingHandler(s.svc)
}

// SetupTest dijalankan sebelum SETIAP test case.
// Bersihkan data test agar setiap test mulai dari kondisi bersih (isolasi).
func (s *TrackingFunctionalSuite) SetupTest() {
	s.awbSeq++
	// Data cleanup akan dibersihkan melalui AWB unik per test,
	// sehingga tidak perlu truncate collection.
}

// generateAWB membuat AWB unik untuk setiap test case
func (s *TrackingFunctionalSuite) generateAWB(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), s.awbSeq)
}

// =========================================================
// FT-01: Test alur lengkap — Catat event → Baca riwayat
// =========================================================

func (s *TrackingFunctionalSuite) TestRecordAndGetHistory() {
	awb := s.generateAWB("FUNC-JKT")

	// --- Step 1: Catat event INBOUND ---
	req1 := &domain.AddTrackingEventRequest{
		AWB:      awb,
		Status:   domain.StatusInbound,
		HubID:    "HUB-JKT-01",
		Location: "Hub Jakarta Pusat",
		Source:   "wms-service",
	}
	event1, err := s.svc.RecordEvent(s.ctx, req1)
	s.NoError(err, "RecordEvent INBOUND harus sukses")
	s.NotNil(event1)
	s.NotEmpty(event1.ID)

	// --- Step 2: Catat event ON_TRANSIT ---
	time.Sleep(10 * time.Millisecond) // Pastikan timestamp berbeda
	req2 := &domain.AddTrackingEventRequest{
		AWB:      awb,
		Status:   domain.StatusOnTransit,
		HubID:    "HUB-JKT-01",
		Location: "Hub Jakarta Pusat",
		Source:   "dispatch-service",
	}
	event2, err := s.svc.RecordEvent(s.ctx, req2)
	s.NoError(err, "RecordEvent ON_TRANSIT harus sukses")
	s.NotNil(event2)

	// --- Step 3: Baca riwayat — harus ada 2 event ---
	history, err := s.svc.GetTrackingHistory(s.ctx, awb)
	s.NoError(err, "GetTrackingHistory harus sukses")
	s.NotNil(history)
	s.Equal(awb, history.AWB)
	s.Equal(2, history.Total)
	s.Len(history.Events, 2)

	// Verifikasi urutan kronologis (INBOUND harus lebih dulu dari ON_TRANSIT)
	s.Equal(domain.StatusInbound, history.Events[0].Status)
	s.Equal(domain.StatusOnTransit, history.Events[1].Status)
	s.True(history.Events[0].Timestamp.Before(history.Events[1].Timestamp) ||
		history.Events[0].Timestamp.Equal(history.Events[1].Timestamp))
}

// =========================================================
// FT-02: Test Cache-Aside Pattern — Redis caching bekerja
// =========================================================

func (s *TrackingFunctionalSuite) TestGetCurrentStatus_CacheAside() {
	awb := s.generateAWB("FUNC-CACHE")

	// --- Step 1: Catat satu event ---
	req := &domain.AddTrackingEventRequest{
		AWB:      awb,
		Status:   domain.StatusAtHub,
		HubID:    "HUB-SBY-01",
		Location: "Hub Surabaya",
		Source:   "wms-service",
	}
	_, err := s.svc.RecordEvent(s.ctx, req)
	s.Require().NoError(err)

	// --- Step 2: GetCurrentStatus — pertama kali (seharusnya dari Redis setelah RecordEvent mengisi cache) ---
	status, err := s.svc.GetCurrentStatus(s.ctx, awb)
	s.NoError(err, "GetCurrentStatus harus sukses")
	s.NotNil(status)
	s.Equal(domain.StatusAtHub, status.CurrentStatus)
	s.Equal("Hub Surabaya", status.LastLocation)

	// --- Step 3: Call kedua — harus dari Redis (faster path) ---
	status2, err := s.svc.GetCurrentStatus(s.ctx, awb)
	s.NoError(err)
	s.Equal(status.CurrentStatus, status2.CurrentStatus)
}

// =========================================================
// FT-03: Test Multi-Event — Status selalu reflect event terbaru
// =========================================================

func (s *TrackingFunctionalSuite) TestStatusAlwaysLatest() {
	awb := s.generateAWB("FUNC-LATEST")

	statuses := []string{
		domain.StatusInbound,
		domain.StatusOnTransit,
		domain.StatusAtHub,
		domain.StatusOutForDelivery,
		domain.StatusDelivered,
	}

	for i, status := range statuses {
		time.Sleep(5 * time.Millisecond) // Pastikan timestamp unik
		req := &domain.AddTrackingEventRequest{
			AWB:      awb,
			Status:   status,
			HubID:    fmt.Sprintf("HUB-%02d", i),
			Location: fmt.Sprintf("Lokasi-%02d", i),
			Source:   "test",
		}
		_, err := s.svc.RecordEvent(s.ctx, req)
		s.Require().NoError(err, "Gagal catat event %s", status)
	}

	// Status terakhir harus DELIVERED
	currentStatus, err := s.svc.GetCurrentStatus(s.ctx, awb)
	s.NoError(err)
	s.Equal(domain.StatusDelivered, currentStatus.CurrentStatus)

	// Total events harus 5
	history, err := s.svc.GetTrackingHistory(s.ctx, awb)
	s.NoError(err)
	s.Equal(5, history.Total)
}

// =========================================================
// FT-04: Test HTTP API end-to-end — Record Event via HTTP
// =========================================================

func (s *TrackingFunctionalSuite) TestHTTPRecordEvent() {
	awb := s.generateAWB("FUNC-HTTP")

	reqBody := domain.AddTrackingEventRequest{
		AWB:      awb,
		Status:   domain.StatusInbound,
		HubID:    "HUB-MKS-01",
		Location: "Hub Makassar",
		Source:   "http-test",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tracking/events",
		strings.NewReader(string(bodyBytes)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.RecordEvent(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)

	var responseEvent domain.TrackingEvent
	err := json.NewDecoder(w.Body).Decode(&responseEvent)
	s.NoError(err)
	s.Equal(awb, responseEvent.AWB)
	s.Equal(domain.StatusInbound, responseEvent.Status)
	s.NotEmpty(responseEvent.ID)
}

// =========================================================
// FT-05: Test HTTP API end-to-end — Get History via HTTP
// =========================================================

func (s *TrackingFunctionalSuite) TestHTTPGetHistory() {
	awb := s.generateAWB("FUNC-HTTP-HIST")

	// Catat beberapa event dulu
	for i, status := range []string{domain.StatusInbound, domain.StatusOnTransit} {
		req := &domain.AddTrackingEventRequest{
			AWB:      awb,
			Status:   status,
			HubID:    fmt.Sprintf("HUB-%02d", i),
			Location: fmt.Sprintf("Lokasi %d", i),
			Source:   "test",
		}
		_, err := s.svc.RecordEvent(s.ctx, req)
		s.Require().NoError(err)
		time.Sleep(5 * time.Millisecond)
	}

	// Panggil HTTP endpoint
	req := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/api/v1/tracking/%s/history", awb), nil)
	w := httptest.NewRecorder()

	s.handler.GetTrackingHistory(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)

	var history domain.TrackingHistory
	err := json.NewDecoder(w.Body).Decode(&history)
	s.NoError(err)
	s.Equal(awb, history.AWB)
	s.Equal(2, history.Total)
}

// =========================================================
// FT-06: Test validasi — Status tidak valid ditolak
// =========================================================

func (s *TrackingFunctionalSuite) TestInvalidStatusRejected() {
	awb := s.generateAWB("FUNC-INVALID")

	req := &domain.AddTrackingEventRequest{
		AWB:    awb,
		Status: "INVALID_STATUS_XYZ",
		HubID:  "HUB-01",
	}

	event, err := s.svc.RecordEvent(s.ctx, req)
	s.Error(err, "Status invalid harus ditolak")
	s.Nil(event)
	s.Contains(err.Error(), "tidak valid")

	// Pastikan tidak ada data tersimpan di MongoDB
	history, err := s.svc.GetTrackingHistory(s.ctx, awb)
	s.NoError(err)
	s.Equal(0, history.Total, "Tidak boleh ada event tersimpan untuk AWB dengan status invalid")
}

// =========================================================
// FT-07: Test AWB tidak ditemukan
// =========================================================

func (s *TrackingFunctionalSuite) TestAWBNotFound() {
	awb := "AWB-YANG-TIDAK-PERNAH-ADA-999999"

	status, err := s.svc.GetCurrentStatus(s.ctx, awb)
	s.Error(err, "Harus error untuk AWB yang tidak ada")
	s.Nil(status)
	s.Contains(err.Error(), "tidak ditemukan")
}

// =========================================================
// NoOpKafkaProducer — Implementasi Kafka yang tidak melakukan apa-apa
// Digunakan di functional test agar tidak butuh Kafka nyata
// =========================================================

type NoOpKafkaProducer struct{}

func (k *NoOpKafkaProducer) PublishEvent(ctx context.Context, topic, key string, value []byte) error {
	return nil // Diabaikan di functional test
}

// =========================================================
// Entry point — Jalankan semua test dalam suite
// =========================================================

func TestTrackingFunctionalSuite(t *testing.T) {
	suite.Run(t, new(TrackingFunctionalSuite))
}