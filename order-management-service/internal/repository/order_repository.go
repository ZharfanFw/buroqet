package repository

import (
	"context"

	"order-management-service/internal/domain"

	"gorm.io/gorm"
)

//go:generate mockgen -source=order_repository.go -destination=../../mock/mock_order_repository.go -package=mock

// OrderRepository defines the contract for order persistence operations.
// This interface is used so the service layer can be tested without a real DB.
type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByAWB(ctx context.Context, awbNumber string) (*model.Order, error)
	FindByTransactionID(ctx context.Context, transactionID string) (*model.Order, error)
	UpdateStatus(ctx context.Context, awbNumber string, status model.OrderStatus) error
}

// orderRepository is the GORM-backed implementation of OrderRepository.
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance backed by the given *gorm.DB.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// Create persists a new order record. It relies on GORM to write all fields
// and will surface any DB-level constraint violations (e.g. duplicate AWB).
func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// FindByAWB retrieves an order by its unique AWB number.
// Returns gorm.ErrRecordNotFound when no matching row exists.
func (r *orderRepository) FindByAWB(ctx context.Context, awbNumber string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Where("awb_number = ?", awbNumber).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByTransactionID retrieves an order by its unique transaction ID.
func (r *orderRepository) FindByTransactionID(ctx context.Context, transactionID string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Where("transaction_id = ?", transactionID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateStatus changes the lifecycle status of an existing order.
func (r *orderRepository) UpdateStatus(ctx context.Context, awbNumber string, status model.OrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.Order{}).
		Where("awb_number = ?", awbNumber).
		Update("status", status).Error
}