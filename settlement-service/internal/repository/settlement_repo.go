package repository

import (
	"context"

	"gorm.io/gorm"
	"settlement-service/internal/domain"
)

type settlementRepository struct {
	db *gorm.DB
}

func NewSettlementRepository(db *gorm.DB) domain.SettlementRepository {
	return &settlementRepository{db: db}
}

func (r *settlementRepository) CreateCommissionLog(ctx context.Context, log *domain.CommissionLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *settlementRepository) GetCommissionsByCourier(ctx context.Context, courierID string) ([]domain.CommissionLog, error) {
	var logs []domain.CommissionLog
	err := r.db.WithContext(ctx).
		Where("courier_id = ?", courierID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *settlementRepository) GetCourierSummary(ctx context.Context, courierID string) (*domain.CourierSummary, error) {
	var logs []domain.CommissionLog
	err := r.db.WithContext(ctx).
		Where("courier_id = ?", courierID).
		Find(&logs).Error
	if err != nil {
		return nil, err
	}

	summary := &domain.CourierSummary{CourierID: courierID}
	for _, log := range logs {
		summary.TotalDeliveries++
		summary.TotalAmount += log.Amount
		if log.Status == "PENDING" {
			summary.PendingAmount += log.Amount
		} else if log.Status == "PAID" {
			summary.PaidAmount += log.Amount
		}
	}
	return summary, nil
}

func (r *settlementRepository) MarkAsPaid(ctx context.Context, courierID string) error {
	return r.db.WithContext(ctx).
		Model(&domain.CommissionLog{}).
		Where("courier_id = ? AND status = ?", courierID, "PENDING").
		Update("status", "PAID").Error
}