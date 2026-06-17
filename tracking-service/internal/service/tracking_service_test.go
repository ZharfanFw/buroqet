// internal/service/tracking_service_test.go
// UNIT TEST untuk TrackingService.
// Semua dependensi (MongoDB, Redis, Kafka) di-MOCK menggunakan gomock.
// Test ini TIDAK membutuhkan database, redis, atau kafka yang berjalan.
// Dapat dijalankan dengan: go test -v -count=1 ./internal/...

package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"tracking-service/internal/domain"
	"tracking-service/internal/service"
	mock_domain "tracking-service/mocks"
)

// setupMocks adalah helper untuk mengurangi boilerplate setup di setiap test case
func setupMocks(t *testing.T) (
	*gomock.Controller,
	*mock_domain.MockTrackingRepository,
	*mock_domain.MockTrackingCache,
	*mock_domain.MockKafkaProducer,
	*service.TrackingService,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mock_domain.NewMockTrackingRepository(ctrl)
	mockCache := mock_domain.NewMockTrackingCache(ctrl)
	mockKafka := mock_domain.NewMockKafkaProducer(ctrl)
	svc := service.NewTrackingService(mockRepo, mockCache, mockKafka)
	return ctrl, mockRepo, mockCache, mockKafka, svc
}

// ============================================================
// TEST SUITE: RecordEvent
// Menguji logika pencatatan event baru ke sistem tracking
// ============================================================

func TestRecordEvent(t *testing.T) {
	ctx := context.Background()

	// ---------------------------------------------------------
	// TC-01: Skenario SUKSES — event baru berhasil dicatat
	// Ekspektasi: InsertEvent dipanggil, Redis diupdate, Kafka dipublish
	// ---------------------------------------------------------
	t.Run("sukses_catat_event_baru", func(t *testing.T) {
		ctrl, mockRepo, mockCache, mockKafka, svc := setupMocks(t)
		defer ctrl.Finish()

		req := &domain.AddTrackingEventRequest{
			AWB:       "JKT-2024-001",
			Status:    domain.StatusOnTransit,
			HubID:     "HUB-JKT-01",
			Location:  "Hub Jakarta Barat",
			Timestamp: time.Now(),
			Source:    "dispatch-service",
		}

		// MongoDB harus menerima insert event
		mockRepo.EXPECT().
			InsertEvent(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, event *domain.TrackingEvent) error {
				// Verifikasi bahwa data di dalam event sudah benar
				assert.Equal(t, req.AWB, event.AWB)
				assert.Equal(t, req.Status, event.Status)
				assert.Equal(t, req.HubID, event.HubID)
				return nil
			}).
			Times(1)

		// Redis harus menerima update status
		mockCache.EXPECT().
			SetStatus(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		// Kafka harus menerima publish event
		mockKafka.EXPECT().
			PublishEvent(gomock.Any(), "tracking.updated", req.AWB, gomock.Any()).
			Return(nil).
			Times(1)

		result, err := svc.RecordEvent(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.AWB, result.AWB)
		assert.Equal(t, domain.StatusOnTransit, result.Status)
		assert.NotEmpty(t, result.ID) // UUID harus di-generate
	})

	// ---------------------------------------------------------
	// TC-02: GAGAL — AWB kosong (validasi input)
	// Ekspektasi: Tidak ada panggilan ke repo/cache/kafka
	// ---------------------------------------------------------
	t.Run("gagal_awb_kosong", func(t *testing.T) {
		ctrl, mockRepo, mockCache, mockKafka, svc := setupMocks(t)
		defer ctrl.Finish()

		// Tidak ada ekspektasi apapun — validasi harus fail sebelum memanggil dependency
		mockRepo.EXPECT().InsertEvent(gomock.Any(), gomock.Any()).Times(0)
		mockCache.EXPECT().SetStatus(gomock.Any(), gomock.Any()).Times(0)
		mockKafka.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		req := &domain.AddTrackingEventRequest{
			AWB:    "", // kosong!
			Status: domain.StatusOnTransit,
			HubID:  "HUB-JKT-01",
		}

		result, err := svc.RecordEvent(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "AWB tidak boleh kosong")
	})

	// ---------------------------------------------------------
	// TC-03: GAGAL — Status tidak valid
	// Ekspektasi: Tidak ada panggilan ke dependency
	// ---------------------------------------------------------
	t.Run("gagal_status_tidak_valid", func(t *testing.T) {
		ctrl, mockRepo, mockCache, mockKafka, svc := setupMocks(t)
		defer ctrl.Finish()

		mockRepo.EXPECT().InsertEvent(gomock.Any(), gomock.Any()).Times(0)
		mockCache.EXPECT().SetStatus(gomock.Any(), gomock.Any()).Times(0)
		mockKafka.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		req := &domain.AddTrackingEventRequest{
			AWB:    "JKT-2024-002",
			Status: "STATUS_TIDAK_DIKENAL", // status invalid!
			HubID:  "HUB-JKT-01",
		}

		result, err := svc.RecordEvent(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tidak valid")
	})

	// ---------------------------------------------------------
	// TC-04: GAGAL — MongoDB error saat InsertEvent
	// Ekspektasi: Error dikembalikan, cache & kafka tidak dipanggil
	// ---------------------------------------------------------
	t.Run("gagal_mongodb_error", func(t *testing.T) {
		ctrl, mockRepo, mockCache, mockKafka, svc := setupMocks(t)
		defer ctrl.Finish()

		req := &domain.AddTrackingEventRequest{
			AWB:    "JKT-2024-003",
			Status: domain.StatusInbound,
			HubID:  "HUB-BDG-01",
		}

		// MongoDB gagal (misal: disk penuh, koneksi terputus)
		mockRepo.EXPECT().
			InsertEvent(gomock.Any(), gomock.Any()).
			Return(errors.New("mongodb: write concern error")).
			Times(1)

		// Cache dan Kafka TIDAK boleh dipanggil kalau DB gagal
		mockCache.EXPECT().SetStatus(gomock.Any(), gomock.Any()).Times(0)
		mockKafka.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		result, err := svc.RecordEvent(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "gagal menyimpan tracking event")
	})

	// ---------------------------------------------------------
	// TC-05: SUKSES meski Redis gagal (eventual consistency)
	// Ekspektasi: Error Redis di-log tapi fungsi tetap return sukses
	// ---------------------------------------------------------
	t.Run("sukses_meski_redis_gagal", func(t *testing.T) {
		ctrl, mockRepo, mockCache, mockKafka, svc := setupMocks(t)
		defer ctrl.Finish()

		req := &domain.AddTrackingEventRequest{
			AWB:    "JKT-2024-004",
			Status: domain.StatusAtHub,
			HubID:  "HUB-SBY-01",
		}

		mockRepo.EXPECT().InsertEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		// Redis gagal — tapi ini tidak boleh menghentikan proses
		mockCache.EXPECT().
			SetStatus(gomock.Any(), gomock.Any()).
			Return(errors.New("redis: connection refused")).
			Times(1)

		// Kafka tetap dipanggil meski Redis gagal
		mockKafka.EXPECT().
			PublishEvent(gomock.Any(), "tracking.updated", req.AWB, gomock.Any()).
			Return(nil).
			Times(1)

		result, err := svc.RecordEvent(ctx, req)

		// Harus sukses meski Redis gagal!
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	// ---------------------------------------------------------
	// TC-06: SUKSES meski Kafka gagal (eventual consistency)
	// Ekspektasi: Data tersimpan di DB & Redis, Kafka error diabaikan
	// ---------------------------------------------------------
	t.Run("sukses_meski_kafka_gagal", func(t *testing.T) {
		ctrl, mockRepo, mockCache, mockKafka, svc := setupMocks(t)
		defer ctrl.Finish()

		req := &domain.AddTrackingEventRequest{
			AWB:    "JKT-2024-005",
			Status: domain.StatusDelivered,
			HubID:  "HUB-DLV-01",
		}

		mockRepo.EXPECT().InsertEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockCache.EXPECT().SetStatus(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		// Kafka broker down
		mockKafka.EXPECT().
			PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errors.New("kafka: leader not available")).
			Times(1)

		result, err := svc.RecordEvent(ctx, req)

		// Harus sukses meski Kafka gagal!
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	// ---------------------------------------------------------
	// TC-07: Timestamp otomatis diisi jika tidak dikirim
	// ---------------------------------------------------------
	t.Run("sukses_timestamp_auto_fill", func(t *testing.T) {
		ctrl, mockRepo, mockCache, mockKafka, svc := setupMocks(t)
		defer ctrl.Finish()

		req := &domain.AddTrackingEventRequest{
			AWB:    "JKT-2024-006",
			Status: domain.StatusInbound,
			HubID:  "HUB-JKT-02",
			// Timestamp sengaja tidak diisi (zero value)
		}

		beforeCall := time.Now()

		mockRepo.EXPECT().
			InsertEvent(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, event *domain.TrackingEvent) error {
				// Timestamp harus sudah terisi otomatis dengan waktu sekarang
				assert.False(t, event.Timestamp.IsZero(), "Timestamp tidak boleh zero")
				assert.True(t, event.Timestamp.After(beforeCall.Add(-time.Second)))
				return nil
			}).
			Times(1)

		mockCache.EXPECT().SetStatus(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockKafka.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		result, err := svc.RecordEvent(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.Timestamp.IsZero())
	})
}

// ============================================================
// TEST SUITE: GetTrackingHistory
// Menguji pembacaan riwayat perjalanan paket
// ============================================================

func TestGetTrackingHistory(t *testing.T) {
	ctx := context.Background()

	// ---------------------------------------------------------
	// TC-08: SUKSES — riwayat ditemukan di MongoDB
	// ---------------------------------------------------------
	t.Run("sukses_dapatkan_riwayat", func(t *testing.T) {
		ctrl, mockRepo, _, _, svc := setupMocks(t)
		defer ctrl.Finish()

		awb := "JKT-2024-007"
		now := time.Now()

		// Simulasikan MongoDB mengembalikan 3 event
		events := []domain.TrackingEvent{
			{ID: "1", AWB: awb, Status: domain.StatusInbound, Timestamp: now.Add(-2 * time.Hour)},
			{ID: "2", AWB: awb, Status: domain.StatusOnTransit, Timestamp: now.Add(-1 * time.Hour)},
			{ID: "3", AWB: awb, Status: domain.StatusAtHub, Timestamp: now},
		}

		mockRepo.EXPECT().
			GetEventsByAWB(gomock.Any(), awb).
			Return(events, nil).
			Times(1)

		result, err := svc.GetTrackingHistory(ctx, awb)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, awb, result.AWB)
		assert.Equal(t, 3, result.Total)
		assert.Len(t, result.Events, 3)
		// Pastikan urutan event benar (ascending by timestamp)
		assert.Equal(t, domain.StatusInbound, result.Events[0].Status)
		assert.Equal(t, domain.StatusAtHub, result.Events[2].Status)
	})

	// ---------------------------------------------------------
	// TC-09: GAGAL — AWB kosong
	// ---------------------------------------------------------
	t.Run("gagal_awb_kosong", func(t *testing.T) {
		ctrl, mockRepo, _, _, svc := setupMocks(t)
		defer ctrl.Finish()

		mockRepo.EXPECT().GetEventsByAWB(gomock.Any(), gomock.Any()).Times(0)

		result, err := svc.GetTrackingHistory(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "AWB tidak boleh kosong")
	})

	// ---------------------------------------------------------
	// TC-10: GAGAL — MongoDB error
	// ---------------------------------------------------------
	t.Run("gagal_mongodb_error", func(t *testing.T) {
		ctrl, mockRepo, _, _, svc := setupMocks(t)
		defer ctrl.Finish()

		awb := "JKT-2024-008"

		mockRepo.EXPECT().
			GetEventsByAWB(gomock.Any(), awb).
			Return(nil, errors.New("mongodb: cursor exhausted")).
			Times(1)

		result, err := svc.GetTrackingHistory(ctx, awb)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "gagal mengambil riwayat tracking")
	})

	// ---------------------------------------------------------
	// TC-11: SUKSES — AWB belum ada event (array kosong)
	// ---------------------------------------------------------
	t.Run("sukses_awb_belum_ada_event", func(t *testing.T) {
		ctrl, mockRepo, _, _, svc := setupMocks(t)
		defer ctrl.Finish()

		awb := "JKT-2024-NEW"

		mockRepo.EXPECT().
			GetEventsByAWB(gomock.Any(), awb).
			Return([]domain.TrackingEvent{}, nil). // array kosong, bukan error
			Times(1)

		result, err := svc.GetTrackingHistory(ctx, awb)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.Total)
		assert.Empty(t, result.Events)
	})
}

// ============================================================
// TEST SUITE: GetCurrentStatus
// Menguji pembacaan status terakhir dengan cache-aside pattern
// ============================================================

func TestGetCurrentStatus(t *testing.T) {
	ctx := context.Background()

	// ---------------------------------------------------------
	// TC-12: SUKSES — Status ditemukan di Redis (cache hit)
	// Ekspektasi: MongoDB TIDAK dipanggil
	// ---------------------------------------------------------
	t.Run("sukses_cache_hit_dari_redis", func(t *testing.T) {
		ctrl, mockRepo, mockCache, _, svc := setupMocks(t)
		defer ctrl.Finish()

		awb := "JKT-2024-009"
		cachedStatus := &domain.TrackingStatus{
			AWB:           awb,
			CurrentStatus: domain.StatusOnTransit,
			LastLocation:  "Hub Surabaya",
			LastUpdated:   time.Now(),
		}

		// Redis berhasil mengembalikan status
		mockCache.EXPECT().
			GetStatus(gomock.Any(), awb).
			Return(cachedStatus, nil).
			Times(1)

		// MongoDB TIDAK boleh dipanggil kalau cache hit
		mockRepo.EXPECT().GetLatestEventByAWB(gomock.Any(), gomock.Any()).Times(0)

		result, err := svc.GetCurrentStatus(ctx, awb)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.StatusOnTransit, result.CurrentStatus)
	})

	// ---------------------------------------------------------
	// TC-13: SUKSES — Cache miss, fallback ke MongoDB
	// Ekspektasi: Redis miss → MongoDB dipanggil → Redis diisi ulang
	// ---------------------------------------------------------
	t.Run("sukses_cache_miss_fallback_mongodb", func(t *testing.T) {
		ctrl, mockRepo, mockCache, _, svc := setupMocks(t)
		defer ctrl.Finish()

		awb := "JKT-2024-010"
		latestEvent := &domain.TrackingEvent{
			ID:        "evt-1",
			AWB:       awb,
			Status:    domain.StatusAtHub,
			Location:  "Hub Bandung",
			Timestamp: time.Now(),
		}

		// Redis tidak punya data (cache miss)
		mockCache.EXPECT().
			GetStatus(gomock.Any(), awb).
			Return(nil, errors.New("redis: key not found")).
			Times(1)

		// MongoDB dipanggil sebagai fallback
		mockRepo.EXPECT().
			GetLatestEventByAWB(gomock.Any(), awb).
			Return(latestEvent, nil).
			Times(1)

		// Redis diisi ulang (cache warming)
		mockCache.EXPECT().
			SetStatus(gomock.Any(), gomock.Any()).
			Return(nil).
			Times(1)

		result, err := svc.GetCurrentStatus(ctx, awb)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, domain.StatusAtHub, result.CurrentStatus)
		assert.Equal(t, "Hub Bandung", result.LastLocation)
	})

	// ---------------------------------------------------------
	// TC-14: GAGAL — AWB kosong
	// ---------------------------------------------------------
	t.Run("gagal_awb_kosong", func(t *testing.T) {
		ctrl, mockRepo, mockCache, _, svc := setupMocks(t)
		defer ctrl.Finish()

		mockCache.EXPECT().GetStatus(gomock.Any(), gomock.Any()).Times(0)
		mockRepo.EXPECT().GetLatestEventByAWB(gomock.Any(), gomock.Any()).Times(0)

		result, err := svc.GetCurrentStatus(ctx, "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "AWB tidak boleh kosong")
	})

	// ---------------------------------------------------------
	// TC-15: GAGAL — Cache miss dan MongoDB juga error
	// ---------------------------------------------------------
	t.Run("gagal_cache_miss_dan_mongodb_error", func(t *testing.T) {
		ctrl, mockRepo, mockCache, _, svc := setupMocks(t)
		defer ctrl.Finish()

		awb := "JKT-2024-011"

		// Cache miss
		mockCache.EXPECT().
			GetStatus(gomock.Any(), awb).
			Return(nil, errors.New("redis: connection refused")).
			Times(1)

		// MongoDB juga error
		mockRepo.EXPECT().
			GetLatestEventByAWB(gomock.Any(), awb).
			Return(nil, errors.New("mongodb: network timeout")).
			Times(1)

		result, err := svc.GetCurrentStatus(ctx, awb)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "gagal mengambil status terkini")
	})

	// ---------------------------------------------------------
	// TC-16: GAGAL — AWB tidak ditemukan sama sekali
	// ---------------------------------------------------------
	t.Run("gagal_awb_tidak_ditemukan", func(t *testing.T) {
		ctrl, mockRepo, mockCache, _, svc := setupMocks(t)
		defer ctrl.Finish()

		awb := "AWB-TIDAK-ADA"

		mockCache.EXPECT().
			GetStatus(gomock.Any(), awb).
			Return(nil, errors.New("key not found")).
			Times(1)

		// MongoDB juga tidak punya data untuk AWB ini
		mockRepo.EXPECT().
			GetLatestEventByAWB(gomock.Any(), awb).
			Return(nil, nil). // nil, nil = tidak ditemukan (bukan error)
			Times(1)

		result, err := svc.GetCurrentStatus(ctx, awb)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "tidak ditemukan")
	})
}

// ============================================================
// TEST SUITE: ProcessKafkaEvent
// Menguji pemrosesan pesan dari Kafka consumer
// ============================================================

func TestProcessKafkaEvent(t *testing.T) {
	ctx := context.Background()

	// ---------------------------------------------------------
	// TC-17: SUKSES — payload Kafka valid berhasil diproses
	// ---------------------------------------------------------
	t.Run("sukses_proses_kafka_event", func(t *testing.T) {
		ctrl, mockRepo, mockCache, mockKafka, svc := setupMocks(t)
		defer ctrl.Finish()

		payload := []byte(`{
			"awb": "JKT-2024-012",
			"status": "ON_TRANSIT",
			"hub_id": "HUB-MKS-01",
			"location": "Hub Makassar",
			"timestamp": "2024-01-15T10:00:00Z",
			"source": "dispatch-service"
		}`)

		mockRepo.EXPECT().InsertEvent(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockCache.EXPECT().SetStatus(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockKafka.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := svc.ProcessKafkaEvent(ctx, payload)

		assert.NoError(t, err)
	})

	// ---------------------------------------------------------
	// TC-18: GAGAL — payload bukan JSON valid
	// ---------------------------------------------------------
	t.Run("gagal_payload_invalid_json", func(t *testing.T) {
		ctrl, mockRepo, mockCache, mockKafka, svc := setupMocks(t)
		defer ctrl.Finish()

		invalidPayload := []byte(`ini bukan JSON`)

		// Tidak ada dependency yang boleh dipanggil
		mockRepo.EXPECT().InsertEvent(gomock.Any(), gomock.Any()).Times(0)
		mockCache.EXPECT().SetStatus(gomock.Any(), gomock.Any()).Times(0)
		mockKafka.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		err := svc.ProcessKafkaEvent(ctx, invalidPayload)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gagal parse kafka payload")
	})
}