// settlement-service/internal/service/settlement_service.go

package service

import (
    "context"
    "fmt"
    "time"

    "settlement-service/internal/domain"
    "github.com/google/uuid"
)

type SettlementService struct {
    repo          domain.SettlementRepository
    pricingClient domain.PricingServiceClient
}

func NewSettlementService(repo domain.SettlementRepository, pricingClient domain.PricingServiceClient) *SettlementService {
    return &SettlementService{repo: repo, pricingClient: pricingClient}
}

// ProcessDeliveryCommission dipanggil ketika event PackageDelivered diterima dari Kafka
func (s *SettlementService) ProcessDeliveryCommission(ctx context.Context, courierID string, awb string, serviceType string) error {
    if courierID == "" || awb == "" {
        return fmt.Errorf("courier ID dan AWB tidak boleh kosong")
    }

    // Ambil commission rate dari Pricing Service (lewat interface = bisa di-mock)
    rate, err := s.pricingClient.GetCommissionRate(ctx, serviceType)
    if err != nil {
        return fmt.Errorf("gagal mengambil commission rate: %w", err)
    }

    if rate <= 0 {
        return fmt.Errorf("commission rate tidak valid: %.2f", rate)
    }

    log := &domain.CommissionLog{
        ID:          uuid.New().String(),
        CourierID:   courierID,
        AWB:         awb,
        Amount:      rate,
        Status:      "PENDING",
        DeliveredAt: time.Now(),
    }

    return s.repo.CreateCommissionLog(ctx, log)
}

// GetCourierEarnings mengambil ringkasan penghasilan kurir
func (s *SettlementService) GetCourierEarnings(ctx context.Context, courierID string) (*domain.CourierSummary, error) {
    if courierID == "" {
        return nil, fmt.Errorf("courier ID tidak boleh kosong")
    }
    return s.repo.GetCourierSummary(ctx, courierID)
}