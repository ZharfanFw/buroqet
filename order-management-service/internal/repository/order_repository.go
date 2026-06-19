package repository

import (
	"context"
	"math"

	"order-management-service/internal/domain"

	"gorm.io/gorm"
)

//go:generate mockgen -source=order_repository.go -destination=../../mock/mock_order_repository.go -package=mock

// OrderRepository defines the contract for order persistence operations.
// This interface is used so the service layer can be tested without a real DB.
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	FindByAWB(ctx context.Context, awbNumber string) (*domain.Order, error)
	FindByTransactionID(ctx context.Context, transactionID string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, awbNumber string, status domain.OrderStatus) error
	FindAll(ctx context.Context, req domain.ListOrdersRequest) (*domain.ListOrdersResponse, error)
	FindByCustomerID(ctx context.Context, customerID string, page, limit int) (*domain.ListOrdersResponse, error)
	DeleteByAWB(ctx context.Context, awbNumber string) error
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
func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// FindByAWB retrieves an order by its unique AWB number.
// Returns gorm.ErrRecordNotFound when no matching row exists.
func (r *orderRepository) FindByAWB(ctx context.Context, awbNumber string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Where("awb_number = ?", awbNumber).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByTransactionID retrieves an order by its payment reference / transaction ID.
func (r *orderRepository) FindByTransactionID(ctx context.Context, transactionID string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Where("payment_ref = ?", transactionID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateStatus changes the lifecycle status of an existing order.
func (r *orderRepository) UpdateStatus(ctx context.Context, awbNumber string, status domain.OrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("awb_number = ?", awbNumber).
		Update("status", status).Error
}

// FindAll retrieves a paginated, optionally filtered list of orders.
// Supports filtering by status and/or customer_id.
func (r *orderRepository) FindAll(ctx context.Context, req domain.ListOrdersRequest) (*domain.ListOrdersResponse, error) {
	// Normalise pagination params
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	query := r.db.WithContext(ctx).Model(&domain.Order{})
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.CustomerID != "" {
		query = query.Where("customer_id = ?", req.CustomerID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var orders []domain.Order
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&orders).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &domain.ListOrdersResponse{
		Orders:     orders,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// FindByCustomerID retrieves a paginated list of orders for a specific customer.
func (r *orderRepository) FindByCustomerID(ctx context.Context, customerID string, page, limit int) (*domain.ListOrdersResponse, error) {
	return r.FindAll(ctx, domain.ListOrdersRequest{
		CustomerID: customerID,
		Page:       page,
		Limit:      limit,
	})
}

// DeleteByAWB permanently removes an order record by AWB number.
func (r *orderRepository) DeleteByAWB(ctx context.Context, awbNumber string) error {
	return r.db.WithContext(ctx).
		Where("awb_number = ?", awbNumber).
		Delete(&domain.Order{}).Error
}

