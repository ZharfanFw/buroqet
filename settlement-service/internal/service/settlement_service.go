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

// ProcessDeliveryCommission dipanggil ketika event PackageDelivered diterima dari Kafka,
// atau saat endpoint POST /api/v1/commissions dipanggil secara manual.
func (s *SettlementService) ProcessDeliveryCommission(ctx context.Context, courierID string, awb string, serviceType string) error {
	if courierID == "" || awb == "" {
		return fmt.Errorf("courier ID dan AWB tidak boleh kosong")
	}

	// Cek apakah AWB sudah pernah diproses sebelumnya (idempotency check)
	existing, err := s.repo.GetCommissionByAWB(ctx, awb)
	if err != nil {
		return fmt.Errorf("gagal memeriksa duplikasi AWB: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("komisi untuk AWB %s sudah pernah dicatat (commission_id: %s)", awb, existing.ID)
	}

	// Ambil commission rate dari Pricing Service (lewat interface = bisa di-mock)
	rate, err := s.pricingClient.GetCommissionRate(ctx, serviceType)
	if err != nil {
		return fmt.Errorf("gagal mengambil commission rate: %w", err)
	}

	if rate <= 0 {
		return fmt.Errorf("commission rate tidak valid: %.2f", rate)
	}

	commissionLog := &domain.CommissionLog{
		ID:          uuid.New().String(),
		CourierID:   courierID,
		AWB:         awb,
		Amount:      rate,
		Status:      "PENDING",
		DeliveredAt: time.Now(),
	}

	return s.repo.CreateCommissionLog(ctx, commissionLog)
}

// GetCourierEarnings mengambil ringkasan penghasilan kurir
func (s *SettlementService) GetCourierEarnings(ctx context.Context, courierID string) (*domain.CourierSummary, error) {
	if courierID == "" {
		return nil, fmt.Errorf("courier ID tidak boleh kosong")
	}
	return s.repo.GetCourierSummary(ctx, courierID)
}

// GetAllCommissions mengambil semua commission log untuk ditampilkan di Frontend.
func (s *SettlementService) GetAllCommissions(ctx context.Context) ([]domain.CommissionLog, error) {
	return s.repo.GetAllCommissions(ctx)
}