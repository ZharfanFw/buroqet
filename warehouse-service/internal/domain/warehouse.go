package domain

import (
    "context"
    "time"
)

// Package merepresentasikan satu paket fisik di gudang
type Package struct {
    ID         string    `json:"id" gorm:"primaryKey"`
    AWB        string    `json:"awb" gorm:"uniqueIndex;not null"`
    HubID      string    `json:"hub_id"`
    ManifestID *string   `json:"manifest_id,omitempty"`
    Status     string    `json:"status"` // INBOUND, OUTBOUND, ON_TRANSIT
    ScannedAt  time.Time `json:"scanned_at"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

// Manifest merepresentasikan satu grouping pengiriman (1 truk/pesawat)
type Manifest struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    TruckID     string    `json:"truck_id"`
    OriginHubID string    `json:"origin_hub_id"`
    DestHubID   string    `json:"dest_hub_id"`
    Status      string    `json:"status"` // OPEN, DISPATCHED
    Packages    []Package `json:"packages" gorm:"foreignKey:ManifestID"`
    CreatedAt   time.Time `json:"created_at"`
}

// WarehouseRepository adalah kontrak yang harus dipenuhi oleh implementasi database
// Interface inilah yang akan di-mock saat unit test
type WarehouseRepository interface {
    SavePackage(ctx context.Context, pkg *Package) error
    GetPackageByAWB(ctx context.Context, awb string) (*Package, error)
    GetAllPackages(ctx context.Context) ([]Package, error)
    UpdatePackageStatus(ctx context.Context, awb string, status string) error
    CreateManifest(ctx context.Context, manifest *Manifest) error
    GetManifestByID(ctx context.Context, id string) (*Manifest, error)
    AddPackageToManifest(ctx context.Context, awb string, manifestID string) error
    DispatchManifest(ctx context.Context, manifestID string) ([]string, error)
}

// KafkaProducer adalah kontrak untuk message broker
// Ini juga perlu di-mock agar unit test tidak butuh Kafka nyata
type KafkaProducer interface {
    PublishEvent(ctx context.Context, topic string, key string, value []byte) error
}