// warehouse-service/internal/service/warehouse_service.go
// Business logic layer untuk Warehouse Management Service.
// Layer ini tidak bergantung pada framework HTTP atau Kafka konkret —
// semua dependency masuk lewat interface (WarehouseRepository, KafkaProducer).

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"warehouse-service/internal/domain"
)

// WarehouseService berisi business logic utama WMS
type WarehouseService struct {
	repo  domain.WarehouseRepository // Tipe interface, bukan konkret!
	kafka domain.KafkaProducer       // Tipe interface, bukan konkret!
}

// NewWarehouseService membuat instance service baru.
// Dengan menerima interface, kita bisa inject mock saat testing.
func NewWarehouseService(repo domain.WarehouseRepository, kafka domain.KafkaProducer) *WarehouseService {
	return &WarehouseService{
		repo:  repo,
		kafka: kafka,
	}
}

// ProcessInbound mencatat paket yang masuk ke gudang.
// Workflow: validasi input → cek duplikat AWB → simpan ke DB → publish package.arrived ke Kafka.
// Side effect: publish event `package.arrived` ke Kafka (sesuai Kafka Event Catalog).
func (s *WarehouseService) ProcessInbound(ctx context.Context, awb string, hubID string) (*domain.Package, error) {
	// Validasi input — logika bisnis
	if awb == "" {
		return nil, fmt.Errorf("AWB tidak boleh kosong")
	}
	if hubID == "" {
		return nil, fmt.Errorf("Hub ID tidak boleh kosong")
	}

	// Cek apakah AWB sudah pernah masuk sebelumnya (idempotency)
	existing, err := s.repo.GetPackageByAWB(ctx, awb)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("paket dengan AWB %s sudah terdaftar di gudang", awb)
	}

	// Buat entitas Package baru
	pkg := &domain.Package{
		ID:        uuid.New().String(),
		AWB:       awb,
		HubID:     hubID,
		Status:    "INBOUND",
		ScannedAt: time.Now(),
	}

	// Simpan ke database (lewat interface)
	if err := s.repo.SavePackage(ctx, pkg); err != nil {
		return nil, fmt.Errorf("gagal menyimpan paket: %w", err)
	}

	// Publish event ke Kafka (lewat interface)
	// Topic: package.arrived (sesuai Kafka Event Catalog di project_context.md section 8)
	// Payload sesuai spec: { "awb": "...", "hub_id": "...", "status": "INBOUND" }
	eventPayload, _ := json.Marshal(map[string]string{
		"awb":    awb,
		"hub_id": hubID,
		"status": "INBOUND",
	})
	if err := s.kafka.PublishEvent(ctx, "package.arrived", awb, eventPayload); err != nil {
		// Log WARNING tapi TIDAK fail — eventual consistency pattern
		// Data sudah masuk DB, event Kafka bisa di-retry nanti
		log.Printf("WARNING: gagal publish event package.arrived untuk AWB %s: %v", awb, err)
	}

	return pkg, nil
}

// ListPackages mengambil semua paket dari database untuk ditampilkan di Frontend.
func (s *WarehouseService) ListPackages(ctx context.Context) ([]domain.Package, error) {
	return s.repo.GetAllPackages(ctx)
}

// DispatchManifest mengirim manifest dan update status semua AWB di dalamnya ke ON_TRANSIT.
// Workflow: validasi → cek manifest exists → cek tidak kosong → dispatch DB → publish manifest.dispatched.
// Side effect: publish event `manifest.dispatched` ke Kafka (sesuai Kafka Event Catalog).
func (s *WarehouseService) DispatchManifest(ctx context.Context, manifestID string) error {
	if manifestID == "" {
		return fmt.Errorf("manifest ID tidak boleh kosong")
	}

	// Validasi: cek manifest ada di database
	manifest, err := s.repo.GetManifestByID(ctx, manifestID)
	if err != nil {
		return fmt.Errorf("manifest dengan ID %s tidak ditemukan: %w", manifestID, err)
	}

	// Validasi: manifest tidak boleh kosong (tidak ada paket)
	if len(manifest.Packages) == 0 {
		return fmt.Errorf("manifest %s tidak memiliki paket, dispatch dibatalkan", manifestID)
	}

	// Dispatch: update status manifest dan semua paket ke ON_TRANSIT
	awbs, err := s.repo.DispatchManifest(ctx, manifestID)
	if err != nil {
		return fmt.Errorf("gagal dispatch manifest: %w", err)
	}

	// Publish satu event yang merepresentasikan semua AWB dalam manifest
	// Topic: manifest.dispatched (sesuai Kafka Event Catalog di project_context.md section 8)
	// Payload: { "manifest_id": "...", "awbs": [...], "status": "ON_TRANSIT" }
	eventPayload, _ := json.Marshal(map[string]interface{}{
		"manifest_id": manifestID,
		"awbs":        awbs,
		"status":      "ON_TRANSIT",
	})
	if err := s.kafka.PublishEvent(ctx, "manifest.dispatched", manifestID, eventPayload); err != nil {
		log.Printf("WARNING: gagal publish event manifest.dispatched untuk manifest %s: %v", manifestID, err)
	}

	return nil
}