package service

import (
	"context"
	"fmt"
	"time"

	"order-management-service/internal/domain"
	"order-management-service/internal/kafka"
	"order-management-service/internal/repository"

	"github.com/google/uuid"
)

//go:generate mockgen -source=order_service.go -destination=../../mock/mock_order_service.go -package=mock

// PricingClient defines the contract for calling the external Pricing & Routing Service.
// Using an interface keeps the service unit-testable without a real HTTP call.
type PricingClient interface {
	GetPrice(ctx context.Context, req domain.PricingRequest) (*domain.PricingResponse, error)
}

// OrderService defines the public business operations for the OMS.
type OrderService interface {
	CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.CreateOrderResponse, error)
	GetOrderByAWB(ctx context.Context, awbNumber string) (*domain.Order, error)
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
func (s *orderService) CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.CreateOrderResponse, error) {
	// --- Step 1: Calculate volumetric & chargeable weight ---
	// Industry standard: L x W x H / 6000
	volumetricWeight := (req.Length * req.Width * req.Height) / 6000
	chargeableWeight := req.WeightActual
	if volumetricWeight > req.WeightActual {
		chargeableWeight = volumetricWeight
	}

	// --- Step 2: Get pricing from Pricing & Routing Service ---
	pricingResp, err := s.pricingClient.GetPrice(ctx, domain.PricingRequest{
		OriginPostal: req.OriginPostal,
		DestPostal:   req.DestPostal,
		Weight:       chargeableWeight,
		Length:       req.Length,
		Width:        req.Width,
		Height:       req.Height,
		ServiceType:  req.ServiceType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pricing: %w", err)
	}

	// --- Step 3: Generate unique identifiers ---
	orderID := uuid.New()
	awbNumber := generateAWB()
	paymentRef := uuid.New().String()

	// --- Step 4: Determine payment URL (only for NON-COD orders) ---
	paymentURL := ""
	if req.PaymentType == domain.PaymentNonCOD {
		// In a real system this would call Payment Gateway (Midtrans/Xendit).
		// Here we construct a placeholder URL; the actual integration is out of scope for OMS.
		paymentURL = fmt.Sprintf("https://pay.example.com/invoice/%s", paymentRef)
	}

	var custID, svcID *uuid.UUID
	if req.CustomerID != "" {
		if parsed, err := uuid.Parse(req.CustomerID); err == nil {
			custID = &parsed
		}
	}
	if req.ServiceID != "" {
		if parsed, err := uuid.Parse(req.ServiceID); err == nil {
			svcID = &parsed
		}
	}

	// --- Step 5: Build the Order entity and persist it ---
	order := &domain.Order{
		OrderID:            orderID,
		AWBNumber:          awbNumber,
		CustomerID:         custID,
		ServiceID:          svcID,
		SenderName:         req.SenderName,
		SenderPhone:        req.SenderPhone,
		SenderAddress:      req.SenderAddress,
		OriginPostalCode:   req.OriginPostal,
		OriginLat:          req.OriginLat,
		OriginLng:          req.OriginLng,
		ReceiverName:       req.ReceiverName,
		ReceiverPhone:      req.ReceiverPhone,
		ReceiverAddress:    req.ReceiverAddress,
		DestPostalCode:     req.DestPostal,
		DestLat:            req.DestLat,
		DestLng:            req.DestLng,
		ActualWeightKg:     req.WeightActual,
		LengthCm:           req.Length,
		WidthCm:            req.Width,
		HeightCm:           req.Height,
		VolumetricWeightKg: volumetricWeight,
		PricingMethod:      "VOLUMETRIC",
		BaseTariff:         pricingResp.BaseFare,
		InsuranceFee:       pricingResp.Insurance,
		Discount:           pricingResp.Discount,
		TotalCost:          pricingResp.TotalPrice,
		UseInsurance:       req.UseInsurance,
		PaymentType:        req.PaymentType,
		IsCod:              req.PaymentType == domain.PaymentCOD,
		RouteID:            nil,
		PaymentRef:         paymentRef,
		Status:             domain.StatusOrderCreated,
		Notes:              "",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if chargeableWeight == req.WeightActual {
		order.PricingMethod = "ACTUAL"
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// --- Step 6: Publish OrderCreated event to Kafka ---
	event := kafka.OrderCreatedEvent{
		AWBNumber:     awbNumber,
		TransactionID: paymentRef,
		SenderName:    req.SenderName,
		SenderAddress: req.SenderAddress,
		OriginCity:    "", // deprecated in Supabase
		ReceiverName:  req.ReceiverName,
		DestCity:      "", // deprecated in Supabase
		ServiceType:   string(req.ServiceType),
		TotalPrice:    pricingResp.TotalPrice,
		CreatedAt:     time.Now(),
	}

	if err := s.kafkaProducer.PublishOrderCreated(ctx, event); err != nil {
		fmt.Printf("[WARN] failed to publish OrderCreated event for AWB %s: %v\n", awbNumber, err)
	}

	return &domain.CreateOrderResponse{
		OrderID:       orderID.String(),
		AWBNumber:     awbNumber,
		PaymentRef:    paymentRef,
		Status:        domain.StatusOrderCreated,
		TotalCost:     pricingResp.TotalPrice,
		PaymentURL:    paymentURL,
	}, nil
}

// GetOrderByAWB fetches order details by AWB number.
func (s *orderService) GetOrderByAWB(ctx context.Context, awbNumber string) (*domain.Order, error) {
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
