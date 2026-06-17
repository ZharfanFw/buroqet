// internal/service/tracking_service.go
// Business logic layer untuk Tracking & Status Service.
// Layer ini hanya bergantung pada interface (domain), bukan implementasi konkret,
// sehingga bisa ditest dengan mock tanpa database/redis/kafka nyata.

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tracking-service/internal/domain"

	"github.com/google/uuid"
)

// TrackingService berisi semua business logic untuk pencatatan dan pembacaan tracking.
// Ia menerima interface (bukan konkret) agar bisa di-inject mock saat unit test.
type TrackingService struct {
	repo  domain.TrackingRepository // MongoDB — source of truth untuk riwayat
	cache domain.TrackingCache      // Redis — cache untuk status terakhir (fast read)
	kafka domain.KafkaProducer      // Kafka — publish ulang event ke downstream jika diperlukan
}

// NewTrackingService membuat instance service dengan Dependency Injection.
// Semua dependensi diinject sebagai interface → mudah di-mock.
func NewTrackingService(
	repo domain.TrackingRepository,
	cache domain.TrackingCache,
	kafka domain.KafkaProducer,
) *TrackingService {
	return &TrackingService{
		repo:  repo,
		cache: cache,
		kafka: kafka,
	}
}

// =========================================================
// RecordEvent — Inti dari service ini
// Mencatat satu event pergerakan paket ke MongoDB dan
// memperbarui status terakhir di Redis.
// Dipanggil baik dari HTTP API maupun dari consumer Kafka.
// =========================================================

func (s *TrackingService) RecordEvent(ctx context.Context, req *domain.AddTrackingEventRequest) (*domain.TrackingEvent, error) {
	// --- 1. VALIDASI INPUT ---
	if req.AWB == "" {
		return nil, fmt.Errorf("AWB tidak boleh kosong")
	}
	if req.Status == "" {
		return nil, fmt.Errorf("status tidak boleh kosong")
	}
	if req.HubID == "" {
		return nil, fmt.Errorf("hub ID tidak boleh kosong")
	}
	if !isValidStatus(req.Status) {
		return nil, fmt.Errorf("status '%s' tidak valid", req.Status)
	}

	// --- 2. BUAT EVENT BARU ---
	// Jika timestamp tidak dikirim (zero value), gunakan waktu sekarang
	eventTime := req.Timestamp
	if eventTime.IsZero() {
		eventTime = time.Now()
	}

	event := &domain.TrackingEvent{
		ID:        uuid.New().String(),
		AWB:       req.AWB,
		Status:    req.Status,
		Location:  req.Location,
		HubID:     req.HubID,
		Timestamp: eventTime,
		CreatedAt: time.Now(),
		Source:    req.Source,
	}

	// --- 3. SIMPAN KE MONGODB (Primary Store) ---
	if err := s.repo.InsertEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("gagal menyimpan tracking event: %w", err)
	}

	// --- 4. UPDATE STATUS TERAKHIR DI REDIS (Cache) ---
	// Redis menyimpan snapshot status terakhir untuk query cepat
	status := &domain.TrackingStatus{
		AWB:           req.AWB,
		CurrentStatus: req.Status,
		LastLocation:  req.Location,
		LastUpdated:   eventTime,
	}
	if err := s.cache.SetStatus(ctx, status); err != nil {
		// Cache miss tidak fatal — eventual consistency
		// Data tetap aman di MongoDB
		fmt.Printf("WARNING: gagal update Redis cache untuk AWB %s: %v\n", req.AWB, err)
	}

	// --- 5. PUBLISH EVENT KE KAFKA (Opsional - untuk downstream service) ---
	eventPayload, _ := json.Marshal(map[string]interface{}{
		"awb":       req.AWB,
		"status":    req.Status,
		"hub_id":    req.HubID,
		"location":  req.Location,
		"timestamp": eventTime,
	})
	if err := s.kafka.PublishEvent(ctx, "tracking.updated", req.AWB, eventPayload); err != nil {
		// Kafka failure tidak fatal — eventual consistency pattern
		fmt.Printf("WARNING: gagal publish tracking event ke Kafka untuk AWB %s: %v\n", req.AWB, err)
	}

	return event, nil
}

// =========================================================
// GetTrackingHistory — Baca riwayat lengkap perjalanan paket
// Selalu membaca dari MongoDB (source of truth)
// =========================================================

func (s *TrackingService) GetTrackingHistory(ctx context.Context, awb string) (*domain.TrackingHistory, error) {
	if awb == "" {
		return nil, fmt.Errorf("AWB tidak boleh kosong")
	}

	events, err := s.repo.GetEventsByAWB(ctx, awb)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil riwayat tracking: %w", err)
	}

	return &domain.TrackingHistory{
		AWB:    awb,
		Events: events,
		Total:  len(events),
	}, nil
}

// =========================================================
// GetCurrentStatus — Baca status TERAKHIR paket (fast path)
// Strategi: coba Redis dulu (cache), kalau miss → MongoDB
// =========================================================

func (s *TrackingService) GetCurrentStatus(ctx context.Context, awb string) (*domain.TrackingStatus, error) {
	if awb == "" {
		return nil, fmt.Errorf("AWB tidak boleh kosong")
	}

	// --- CACHE-ASIDE PATTERN ---
	// 1. Coba baca dari Redis (fast)
	cached, err := s.cache.GetStatus(ctx, awb)
	if err == nil && cached != nil {
		// Cache hit → langsung return
		return cached, nil
	}

	// 2. Cache miss → baca dari MongoDB (source of truth)
	latestEvent, err := s.repo.GetLatestEventByAWB(ctx, awb)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil status terkini: %w", err)
	}
	if latestEvent == nil {
		return nil, fmt.Errorf("AWB %s tidak ditemukan", awb)
	}

	// 3. Reconstruct status dari event terbaru
	status := &domain.TrackingStatus{
		AWB:           awb,
		CurrentStatus: latestEvent.Status,
		LastLocation:  latestEvent.Location,
		LastUpdated:   latestEvent.Timestamp,
	}

	// 4. Isi ulang cache (cache warming) agar next request cepat
	if cacheErr := s.cache.SetStatus(ctx, status); cacheErr != nil {
		fmt.Printf("WARNING: gagal warm Redis cache untuk AWB %s: %v\n", awb, cacheErr)
	}

	return status, nil
}

// =========================================================
// ProcessKafkaEvent — Handler untuk event dari Kafka
// Dipanggil oleh Kafka consumer loop ketika ada pesan masuk
// =========================================================

func (s *TrackingService) ProcessKafkaEvent(ctx context.Context, payload []byte) error {
	var kafkaPayload domain.KafkaTrackingPayload
	if err := json.Unmarshal(payload, &kafkaPayload); err != nil {
		return fmt.Errorf("gagal parse kafka payload: %w", err)
	}

	req := &domain.AddTrackingEventRequest{
		AWB:       kafkaPayload.AWB,
		Status:    kafkaPayload.Status,
		HubID:     kafkaPayload.HubID,
		Location:  kafkaPayload.Location,
		Timestamp: kafkaPayload.Timestamp,
		Source:    kafkaPayload.Source,
	}

	_, err := s.RecordEvent(ctx, req)
	return err
}

// =========================================================
// HELPER FUNCTIONS
// =========================================================

// isValidStatus memvalidasi apakah status yang diberikan adalah salah satu
// nilai yang dikenali oleh sistem. Ini adalah business rule.
func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		domain.StatusInbound:         true,
		domain.StatusOnTransit:       true,
		domain.StatusAtHub:           true,
		domain.StatusOutForDelivery:  true,
		domain.StatusDelivered:       true,
		domain.StatusFailed:          true,
		domain.StatusReturned:        true,
	}
	return validStatuses[status]
}