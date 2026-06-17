package service

import (
	"context"
	"fmt"
	"time"

	"order-management-service/internal/kafka"
	"order-management-service/internal/domain"
	"order-management-service/internal/repository"

	"github.com/google/uuid"
)

//go:generate mockgen -source=order_service.go -destination=../../mock/mock_order_service.go -package=mock

// PricingClient defines the contract for calling the external Pricing & Routing Service.
// Using an interface keeps the service unit-testable without a real HTTP call.
type PricingClient interface {
	GetPrice(ctx context.Context, req model.PricingRequest) (*model.PricingResponse, error)
}

// OrderService defines the public business operations for the OMS.
type OrderService interface {
	CreateOrder(ctx context.Context, req model.CreateOrderRequest) (*model.CreateOrderResponse, error)
	GetOrderByAWB(ctx context.Context, awbNumber string) (*model.Order, error)
}

// orderService is the concrete implementation.
type orderService struct {
	repo          repository.OrderRepository
	pricingClient PricingClient
	kafkaProducer kafka.Producer
}

// NewOrderService wires together the dependencies and returns an OrderService.
func NewOrderService(
	repo repository.OrderRepository,
	pricingClient PricingClient,
	kafkaProducer kafka.Producer,
) OrderService {
	return &orderService{
		repo:          repo,
		pricingClient: pricingClient,
		kafkaProducer: kafkaProducer,
	}
}

// CreateOrder orchestrates the full order-creation flow:
//  1. Call Pricing Service (synchronous REST) to get the shipping cost.
//  2. Calculate volumetric weight.
//  3. Generate unique AWB and Transaction ID.
//  4. Persist the order to PostgreSQL.
//  5. Publish OrderCreated event to Kafka for Dispatch Service.
//  6. Return the response DTO to the handler.
func (s *orderService) CreateOrder(ctx context.Context, req model.CreateOrderRequest) (*model.CreateOrderResponse, error) {
	// --- Step 1: Get pricing from Pricing & Routing Service ---
	pricingResp, err := s.pricingClient.GetPrice(ctx, model.PricingRequest{
		OriginPostal: req.OriginPostal,
		DestPostal:   req.DestPostal,
		Weight:       req.WeightActual,
		Length:       req.Length,
		Width:        req.Width,
		Height:       req.Height,
		ServiceType:  req.ServiceType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pricing: %w", err)
	}

	// --- Step 2: Calculate volumetric weight (L x W x H / 6000 is the industry standard) ---
	volumetricWeight := (req.Length * req.Width * req.Height) / 6000

	// --- Step 3: Generate unique identifiers ---
	awbNumber := generateAWB()
	transactionID := uuid.New().String()

	// --- Step 4: Determine payment URL (only for NON-COD orders) ---
	paymentURL := ""
	if req.PaymentType == model.PaymentNonCOD {
		// In a real system this would call Payment Gateway (Midtrans/Xendit).
		// Here we construct a placeholder URL; the actual integration is out of scope for OMS.
		paymentURL = fmt.Sprintf("https://pay.example.com/invoice/%s", transactionID)
	}

	// --- Step 5: Build the Order entity and persist it ---
	order := &model.Order{
		AWBNumber:       awbNumber,
		TransactionID:   transactionID,
		Status:          model.StatusOrderCreated,
		SenderName:      req.SenderName,
		SenderPhone:     req.SenderPhone,
		SenderAddress:   req.SenderAddress,
		OriginCity:      req.OriginCity,
		OriginPostal:    req.OriginPostal,
		ReceiverName:    req.ReceiverName,
		ReceiverPhone:   req.ReceiverPhone,
		ReceiverAddress: req.ReceiverAddress,
		DestCity:        req.DestCity,
		DestPostal:      req.DestPostal,
		WeightActual:    req.WeightActual,
		WeightVolumetri: volumetricWeight,
		Length:          req.Length,
		Width:           req.Width,
		Height:          req.Height,
		ServiceType:     req.ServiceType,
		PaymentType:     req.PaymentType,
		TotalPrice:      pricingResp.TotalPrice,
		PaymentURL:      paymentURL,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// --- Step 6: Publish OrderCreated event to Kafka ---
	event := kafka.OrderCreatedEvent{
		AWBNumber:     awbNumber,
		TransactionID: transactionID,
		SenderName:    req.SenderName,
		SenderAddress: req.SenderAddress,
		OriginCity:    req.OriginCity,
		ReceiverName:  req.ReceiverName,
		DestCity:      req.DestCity,
		ServiceType:   string(req.ServiceType),
		TotalPrice:    pricingResp.TotalPrice,
		CreatedAt:     time.Now(),
	}

	if err := s.kafkaProducer.PublishOrderCreated(ctx, event); err != nil {
		// Non-fatal: order is already persisted. Log and continue.
		// In production you'd push to a dead-letter queue or an outbox table.
		fmt.Printf("[WARN] failed to publish OrderCreated event for AWB %s: %v\n", awbNumber, err)
	}

	return &model.CreateOrderResponse{
		AWBNumber:     awbNumber,
		TransactionID: transactionID,
		Status:        model.StatusOrderCreated,
		TotalPrice:    pricingResp.TotalPrice,
		PaymentURL:    paymentURL,
	}, nil
}

// GetOrderByAWB fetches order details by AWB number.
func (s *orderService) GetOrderByAWB(ctx context.Context, awbNumber string) (*model.Order, error) {
	order, err := s.repo.FindByAWB(ctx, awbNumber)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	return order, nil
}

// generateAWB creates a unique AWB in the format "JNE-<8-char-UUID-prefix>".
// A real implementation would use a more structured format (e.g. branch code + date + sequence).
func generateAWB() string {
	id := uuid.New().String()
	return fmt.Sprintf("JNE-%s", id[:8])
}