// internal/repository/tracking_repo_test.go
// Functional test untuk MongoTrackingRepository.
// Test ini MENGAKSES MongoDB nyata — dijalankan dengan build tag functional:
//   go test -v -tags=functional -count=1 ./internal/repository/...
//
// Tujuan: memverifikasi bahwa implementasi MongoDB (InsertEvent, GetEventsByAWB,
// GetLatestEventByAWB) bekerja benar terhadap database nyata.

//go:build functional
// +build functional

package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"tracking-service/internal/domain"
	"tracking-service/internal/repository"
)


// =========================================================
// MongoRepoSuite — Test suite untuk MongoDB repository
// =========================================================

type MongoRepoSuite struct {
	suite.Suite
	repo *repository.MongoTrackingRepository
	ctx  context.Context
	seq  int
}

func (s *MongoRepoSuite) SetupSuite() {
	s.ctx = context.Background()

	mongoURI := os.Getenv("TEST_MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://testuser:testpass@localhost:27018/tracking_test?authSource=admin"
	}
	mongoDBName := os.Getenv("TEST_MONGO_DB")
	if mongoDBName == "" {
		mongoDBName = "tracking_test"
	}

	db, err := repository.ConnectMongoDB(mongoURI, mongoDBName)
	s.Require().NoError(err, "Gagal koneksi ke MongoDB test")

	s.repo = repository.NewMongoTrackingRepository(db)

	// Buat index
	err = s.repo.EnsureIndexes(s.ctx)
	s.Require().NoError(err, "Gagal buat MongoDB index")
}

func (s *MongoRepoSuite) SetupTest() {
	s.seq++
}

// generateAWB membuat AWB unik per test case
func (s *MongoRepoSuite) generateAWB(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), s.seq)
}

// =========================================================
// REPO-01: InsertEvent — event tersimpan ke MongoDB
// =========================================================

func (s *MongoRepoSuite) TestInsertEvent_Sukses() {
	awb := s.generateAWB("REPO-INSERT")
	event := &domain.TrackingEvent{
		ID:        fmt.Sprintf("id-%d", time.Now().UnixNano()),
		AWB:       awb,
		Status:    domain.StatusInbound,
		Location:  "Hub Jakarta",
		HubID:     "HUB-JKT-01",
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
		Source:    "test",
	}

	err := s.repo.InsertEvent(s.ctx, event)

	// Akan FAIL karena implementasi masih "not implemented"
	// Setelah diimplementasi, harus:
	s.NoError(err, "InsertEvent harus sukses")
}

// =========================================================
// REPO-02: InsertEvent — field ID wajib diisi
// =========================================================

func (s *MongoRepoSuite) TestInsertEvent_TanpaID() {
	awb := s.generateAWB("REPO-NOID")
	event := &domain.TrackingEvent{
		// ID sengaja dikosongkan
		AWB:       awb,
		Status:    domain.StatusInbound,
		Timestamp: time.Now(),
	}

	// MongoDB akan generate _id otomatis jika kosong
	// Implementasi harus handle ini
	err := s.repo.InsertEvent(s.ctx, event)
	s.NoError(err, "InsertEvent tanpa ID tetap harus sukses — MongoDB generate _id otomatis")
}

// =========================================================
// REPO-03: GetEventsByAWB — ambil semua event, urut ascending
// =========================================================

func (s *MongoRepoSuite) TestGetEventsByAWB_Urutan() {
	awb := s.generateAWB("REPO-GETALL")
	now := time.Now()

	// Insert 3 event dengan timestamp berbeda
	events := []*domain.TrackingEvent{
		{ID: "r1", AWB: awb, Status: domain.StatusInbound, Timestamp: now.Add(-2 * time.Hour), CreatedAt: now},
		{ID: "r2", AWB: awb, Status: domain.StatusOnTransit, Timestamp: now.Add(-1 * time.Hour), CreatedAt: now},
		{ID: "r3", AWB: awb, Status: domain.StatusAtHub, Timestamp: now, CreatedAt: now},
	}
	for _, e := range events {
		s.repo.InsertEvent(s.ctx, e)
	}

	result, err := s.repo.GetEventsByAWB(s.ctx, awb)

	// Akan FAIL karena "not implemented"
	s.NoError(err)
	s.Len(result, 3, "Harus mengembalikan 3 event")

	// Verifikasi urutan chronologis (ascending timestamp)
	s.Equal(domain.StatusInbound, result[0].Status,   "Event pertama harus INBOUND")
	s.Equal(domain.StatusOnTransit, result[1].Status, "Event kedua harus ON_TRANSIT")
	s.Equal(domain.StatusAtHub, result[2].Status,     "Event ketiga harus AT_HUB")

	// Pastikan timestamp benar-benar ascending
	s.True(result[0].Timestamp.Before(result[1].Timestamp))
	s.True(result[1].Timestamp.Before(result[2].Timestamp))
}

// =========================================================
// REPO-04: GetEventsByAWB — AWB tidak ditemukan → slice kosong
// =========================================================

func (s *MongoRepoSuite) TestGetEventsByAWB_TidakDitemukan() {
	awb := s.generateAWB("REPO-NOTFOUND")

	result, err := s.repo.GetEventsByAWB(s.ctx, awb)

	s.NoError(err, "AWB tidak ada bukan error, return slice kosong")
	s.NotNil(result)
	s.Empty(result, "Harus return slice kosong, bukan nil")
}

// =========================================================
// REPO-05: GetEventsByAWB — isolasi antar AWB
// Memastikan query tidak mengambil event milik AWB lain
// =========================================================

func (s *MongoRepoSuite) TestGetEventsByAWB_IsolasiAWB() {
	awbA := s.generateAWB("REPO-AWB-A")
	awbB := s.generateAWB("REPO-AWB-B")
	now := time.Now()

	// Insert event untuk dua AWB berbeda
	s.repo.InsertEvent(s.ctx, &domain.TrackingEvent{ID: "iso-a1", AWB: awbA, Status: domain.StatusInbound, Timestamp: now})
	s.repo.InsertEvent(s.ctx, &domain.TrackingEvent{ID: "iso-a2", AWB: awbA, Status: domain.StatusOnTransit, Timestamp: now.Add(time.Minute)})
	s.repo.InsertEvent(s.ctx, &domain.TrackingEvent{ID: "iso-b1", AWB: awbB, Status: domain.StatusInbound, Timestamp: now})

	// Query AWB-A — harus dapat 2 event, bukan 3
	resultA, err := s.repo.GetEventsByAWB(s.ctx, awbA)
	s.NoError(err)
	s.Len(resultA, 2, "AWB-A harus punya 2 event")
	for _, e := range resultA {
		s.Equal(awbA, e.AWB, "Semua event harus milik AWB-A")
	}

	// Query AWB-B — harus dapat 1 event
	resultB, err := s.repo.GetEventsByAWB(s.ctx, awbB)
	s.NoError(err)
	s.Len(resultB, 1, "AWB-B harus punya 1 event")
}

// =========================================================
// REPO-06: GetLatestEventByAWB — ambil event paling baru
// =========================================================

func (s *MongoRepoSuite) TestGetLatestEventByAWB_AmbilPalingBaru() {
	awb := s.generateAWB("REPO-LATEST")
	now := time.Now()

	// Insert beberapa event, yang terakhir adalah DELIVERED
	s.repo.InsertEvent(s.ctx, &domain.TrackingEvent{ID: "l1", AWB: awb, Status: domain.StatusInbound, Timestamp: now.Add(-3 * time.Hour)})
	s.repo.InsertEvent(s.ctx, &domain.TrackingEvent{ID: "l2", AWB: awb, Status: domain.StatusOnTransit, Timestamp: now.Add(-2 * time.Hour)})
	s.repo.InsertEvent(s.ctx, &domain.TrackingEvent{ID: "l3", AWB: awb, Status: domain.StatusDelivered, Timestamp: now.Add(-1 * time.Hour)})

	latest, err := s.repo.GetLatestEventByAWB(s.ctx, awb)

	// Akan FAIL karena "not implemented"
	s.NoError(err)
	s.NotNil(latest)
	// Harus mengembalikan event dengan timestamp paling baru = DELIVERED
	s.Equal(domain.StatusDelivered, latest.Status, "Event terbaru harus DELIVERED")
}

// =========================================================
// REPO-07: GetLatestEventByAWB — AWB tidak ada → return nil, nil
// =========================================================

func (s *MongoRepoSuite) TestGetLatestEventByAWB_TidakDitemukan() {
	awb := s.generateAWB("REPO-LATEST-NOTFOUND")

	latest, err := s.repo.GetLatestEventByAWB(s.ctx, awb)

	s.NoError(err, "AWB tidak ada bukan error")
	s.Nil(latest, "Harus return nil jika AWB tidak ditemukan")
}

// =========================================================
// REPO-08: EnsureIndexes — idempoten (aman dipanggil berulang)
// =========================================================

func (s *MongoRepoSuite) TestEnsureIndexes_Idempoten() {
	// Panggil EnsureIndexes dua kali — tidak boleh error
	err1 := s.repo.EnsureIndexes(s.ctx)
	err2 := s.repo.EnsureIndexes(s.ctx)

	s.NoError(err1, "EnsureIndexes pertama harus sukses")
	s.NoError(err2, "EnsureIndexes kedua (idempoten) juga harus sukses")
}

// =========================================================
// Entry point
// =========================================================

func TestMongoRepoSuite(t *testing.T) {
	suite.Run(t, new(MongoRepoSuite))
}