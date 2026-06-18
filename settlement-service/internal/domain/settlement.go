package domain

import (
	"context"
	"time"
)

// CommissionLog merepresentasikan satu catatan komisi kurir untuk satu pengiriman.
// Disimpan di tabel commission_logs di PostgreSQL (settlement_db).
type CommissionLog struct {
	ID          string    `gorm:"primaryKey"          json:"id"`
	CourierID   string    `gorm:"not null;index"      json:"courier_id"`
	AWB         string    `gorm:"not null;uniqueIndex" json:"awb"` // uniqueIndex untuk mencegah duplikat
	Amount      float64   `gorm:"not null"            json:"amount"`
	Status      string    `gorm:"not null;default:'PENDING'" json:"status"` // PENDING | PAID
	DeliveredAt time.Time `json:"delivered_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CourierSummary adalah DTO agregat untuk rekap penghasilan kurir.
// Tidak disimpan ke DB — dihitung on-the-fly dari commission_logs.
type CourierSummary struct {
	CourierID       string  `json:"courier_id"`
	TotalDeliveries int     `json:"total_deliveries"`
	TotalAmount     float64 `json:"total_amount"`
	PendingAmount   float64 `json:"pending_amount"`
	PaidAmount      float64 `json:"paid_amount"`
}

// SettlementRepository adalah kontrak yang harus dipenuhi oleh implementasi database.
// Interface inilah yang di-mock saat unit test.
type SettlementRepository interface {
	CreateCommissionLog(ctx context.Context, log *CommissionLog) error
	GetCommissionByAWB(ctx context.Context, awb string) (*CommissionLog, error)
	GetCommissionsByCourier(ctx context.Context, courierID string) ([]CommissionLog, error)
	GetCourierSummary(ctx context.Context, courierID string) (*CourierSummary, error)
	MarkAsPaid(ctx context.Context, courierID string) error
}

// PricingServiceClient adalah kontrak untuk memanggil Pricing Service via HTTP.
// Interface ini memungkinkan mock di unit test tanpa butuh HTTP call nyata.
type PricingServiceClient interface {
	GetCommissionRate(ctx context.Context, serviceType string) (float64, error)
}