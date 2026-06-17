package service

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "warehouse-service/internal/domain"
    "github.com/google/uuid"
)

// WarehouseService berisi business logic utama WMS
type WarehouseService struct {
    repo     domain.WarehouseRepository  // Tipe interface, bukan konkret!
    kafka    domain.KafkaProducer        // Tipe interface, bukan konkret!
}

// NewWarehouseService membuat instance service baru
// Dengan menerima interface, kita bisa inject mock saat testing
func NewWarehouseService(repo domain.WarehouseRepository, kafka domain.KafkaProducer) *WarehouseService {
    return &WarehouseService{
        repo:  repo,
        kafka: kafka,
    }
}

// ProcessInbound mencatat paket yang masuk ke gudang
// Ini adalah business logic yang akan kita test
func (s *WarehouseService) ProcessInbound(ctx context.Context, awb string, hubID string) (*domain.Package, error) {
    // Validasi input — logika bisnis sederhana
    if awb == "" {
        return nil, fmt.Errorf("AWB tidak boleh kosong")
    }
    if hubID == "" {
        return nil, fmt.Errorf("Hub ID tidak boleh kosong")
    }

    // Cek apakah AWB sudah pernah masuk sebelumnya
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
    eventPayload, _ := json.Marshal(map[string]string{
        "awb":    awb,
        "hub_id": hubID,
        "status": "INBOUND",
    })
    if err := s.kafka.PublishEvent(ctx, "package.arrived", awb, eventPayload); err != nil {
        // Kita log tapi tidak fail — eventual consistency pattern
        // Artinya data sudah masuk DB, event Kafka bisa di-retry nanti
        fmt.Printf("WARNING: gagal publish event untuk AWB %s: %v\n", awb, err)
    }

    return pkg, nil
}

// DispatchManifest mengirim manifest dan update status semua AWB di dalamnya
func (s *WarehouseService) DispatchManifest(ctx context.Context, manifestID string) error {
    if manifestID == "" {
        return fmt.Errorf("manifest ID tidak boleh kosong")
    }

    // Ambil semua AWB dalam manifest, lalu ubah statusnya ke ON_TRANSIT
    awbs, err := s.repo.DispatchManifest(ctx, manifestID)
    if err != nil {
        return fmt.Errorf("gagal dispatch manifest: %w", err)
    }

    // Publish satu event yang merepresentasikan semua AWB
    eventPayload, _ := json.Marshal(map[string]interface{}{
        "manifest_id": manifestID,
        "awbs":        awbs,
        "status":      "ON_TRANSIT",
    })
    if err := s.kafka.PublishEvent(ctx, "manifest.dispatched", manifestID, eventPayload); err != nil {
        fmt.Printf("WARNING: gagal publish manifest event: %v\n", err)
    }

    return nil
}