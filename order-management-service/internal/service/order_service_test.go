package service_test

import (
	"context"
	"errors"
	"testing"

	"order-management-service/internal/kafka"
	"order-management-service/internal/domain"
	"order-management-service/internal/service"
	"order-management-service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// helper: returns a valid CreateOrderRequest to avoid repetition in every test
func validCreateOrderRequest() model.CreateOrderRequest {
	return model.CreateOrderRequest{
		SenderName:      "Budi Santoso",
		SenderPhone:     "081234567890",
		SenderAddress:   "Jl. Merdeka No.1, Jakarta",
		OriginCity:      "Jakarta",
		OriginPostal:    "10110",
		ReceiverName:    "Ani Rahayu",
		ReceiverPhone:   "089876543210",
		ReceiverAddress: "Jl. Sudirman No.5, Bandung",
		DestCity:        "Bandung",
		DestPostal:      "40111",
		WeightActual:    2.0,
		Length:          20,
		Width:           15,
		Height:          10,
		ServiceType:     model.ServiceRegular,
		PaymentType:     model.PaymentNonCOD,
	}
}

// helper: returns a canned pricing response from the mock pricing client
func validPricingResponse() *model.PricingResponse {
	return &model.PricingResponse{
		BaseFare:     15000,
		Insurance:    2000,
		Discount:     0,
		TotalPrice:   17000,
		EstimatedSLA: "2-3 hari",
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// CreateOrder Tests
// ─────────────────────────────────────────────────────────────────────────────

// TestCreateOrder_Success verifies the happy path:
// pricing client returns a price → order is saved → event is published → response is correct.
func TestCreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)

	req := validCreateOrderRequest()
	ctx := context.Background()

	// Pricing service will be called once with the correct postal codes
	mockPricing.EXPECT().
		GetPrice(ctx, gomock.Any()).
		Return(validPricingResponse(), nil)

	// Repository will be called once to persist the order
	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	// Kafka will be called once to publish the event
	mockKafka.EXPECT().
		PublishOrderCreated(ctx, gomock.Any()).
		Return(nil)

	resp, err := svc.CreateOrder(ctx, req)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.AWBNumber, "AWB number should be generated")
	assert.NotEmpty(t, resp.TransactionID, "Transaction ID should be generated")
	assert.Equal(t, model.StatusOrderCreated, resp.Status)
	assert.Equal(t, 17000.0, resp.TotalPrice)
	assert.NotEmpty(t, resp.PaymentURL, "Non-COD order must have a payment URL")
}

// TestCreateOrder_COD_NoPaymentURL verifies that COD orders do NOT get a payment URL.
func TestCreateOrder_COD_NoPaymentURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)

	req := validCreateOrderRequest()
	req.PaymentType = model.PaymentCOD
	ctx := context.Background()

	mockPricing.EXPECT().GetPrice(ctx, gomock.Any()).Return(validPricingResponse(), nil)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	mockKafka.EXPECT().PublishOrderCreated(ctx, gomock.Any()).Return(nil)

	resp, err := svc.CreateOrder(ctx, req)

	require.NoError(t, err)
	assert.Empty(t, resp.PaymentURL, "COD orders must NOT have a payment URL")
}

// TestCreateOrder_ExpressService verifies that Express orders trigger the correct pricing call
// and the response carries the correct (higher) price.
func TestCreateOrder_ExpressService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)

	req := validCreateOrderRequest()
	req.ServiceType = model.ServiceExpress
	ctx := context.Background()

	expressPricing := &model.PricingResponse{
		BaseFare:     30000,
		Insurance:    2000,
		Discount:     0,
		TotalPrice:   32000,
		EstimatedSLA: "next-day",
	}

	mockPricing.EXPECT().GetPrice(ctx, gomock.Any()).Return(expressPricing, nil)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	mockKafka.EXPECT().PublishOrderCreated(ctx, gomock.Any()).Return(nil)

	resp, err := svc.CreateOrder(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 32000.0, resp.TotalPrice)
}

// TestCreateOrder_PricingServiceError verifies that if the Pricing Service fails,
// CreateOrder returns an error and does NOT touch the database or Kafka.
func TestCreateOrder_PricingServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)

	req := validCreateOrderRequest()
	ctx := context.Background()

	mockPricing.EXPECT().
		GetPrice(ctx, gomock.Any()).
		Return(nil, errors.New("pricing service timeout"))

	// repo and kafka must NOT be called when pricing fails
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
	mockKafka.EXPECT().PublishOrderCreated(gomock.Any(), gomock.Any()).Times(0)

	resp, err := svc.CreateOrder(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get pricing")
}

// TestCreateOrder_RepositoryError verifies that a DB failure after a successful
// pricing call returns an error and does NOT publish to Kafka.
func TestCreateOrder_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)

	ctx := context.Background()
	req := validCreateOrderRequest()

	mockPricing.EXPECT().GetPrice(ctx, gomock.Any()).Return(validPricingResponse(), nil)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("duplicate awb_number"))

	// Kafka must NOT be called if the DB write fails
	mockKafka.EXPECT().PublishOrderCreated(gomock.Any(), gomock.Any()).Times(0)

	resp, err := svc.CreateOrder(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save order")
}

// TestCreateOrder_KafkaError verifies that a Kafka failure is non-fatal:
// the order is still created and the response is returned without error.
func TestCreateOrder_KafkaError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)

	ctx := context.Background()
	req := validCreateOrderRequest()

	mockPricing.EXPECT().GetPrice(ctx, gomock.Any()).Return(validPricingResponse(), nil)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	mockKafka.EXPECT().
		PublishOrderCreated(ctx, gomock.Any()).
		Return(errors.New("kafka broker unavailable"))

	// Even with Kafka failure, the response should still be returned
	resp, err := svc.CreateOrder(ctx, req)

	require.NoError(t, err, "Kafka failure must be non-fatal — order is already persisted")
	assert.NotEmpty(t, resp.AWBNumber)
}

// TestCreateOrder_AWBIsUnique checks that each call generates a different AWB.
func TestCreateOrder_AWBIsUnique(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)
	ctx := context.Background()

	mockPricing.EXPECT().GetPrice(ctx, gomock.Any()).Return(validPricingResponse(), nil).Times(2)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil).Times(2)
	mockKafka.EXPECT().PublishOrderCreated(ctx, gomock.Any()).Return(nil).Times(2)

	resp1, _ := svc.CreateOrder(ctx, validCreateOrderRequest())
	resp2, _ := svc.CreateOrder(ctx, validCreateOrderRequest())

	assert.NotEqual(t, resp1.AWBNumber, resp2.AWBNumber, "Every AWB must be unique")
	assert.NotEqual(t, resp1.TransactionID, resp2.TransactionID)
}

// TestCreateOrder_VolumetricWeightCalculation verifies that volumetric weight
// is calculated correctly (L x W x H / 6000) and stored in the persisted entity.
func TestCreateOrder_VolumetricWeightCalculation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)
	ctx := context.Background()

	req := validCreateOrderRequest()
	// 30 x 20 x 10 / 6000 = 1.0 kg volumetric
	req.Length = 30
	req.Width = 20
	req.Height = 10

	expectedVolumetric := (30.0 * 20.0 * 10.0) / 6000.0 // = 1.0

	mockPricing.EXPECT().GetPrice(ctx, gomock.Any()).Return(validPricingResponse(), nil)

	// Capture the order passed to Create and assert its volumetric weight
	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, order *model.Order) error {
			assert.InDelta(t, expectedVolumetric, order.WeightVolumetri, 0.001,
				"Volumetric weight must be L*W*H/6000")
			return nil
		})

	mockKafka.EXPECT().PublishOrderCreated(ctx, gomock.Any()).Return(nil)

	_, err := svc.CreateOrder(ctx, req)
	require.NoError(t, err)
}

// TestCreateOrder_KafkaEventContainsCorrectAWB verifies that the event published
// to Kafka carries the same AWB as the response.
func TestCreateOrder_KafkaEventContainsCorrectAWB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)
	ctx := context.Background()

	mockPricing.EXPECT().GetPrice(ctx, gomock.Any()).Return(validPricingResponse(), nil)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	var capturedEvent kafka.OrderCreatedEvent
	mockKafka.EXPECT().
		PublishOrderCreated(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, event kafka.OrderCreatedEvent) error {
			capturedEvent = event
			return nil
		})

	resp, err := svc.CreateOrder(ctx, validCreateOrderRequest())
	require.NoError(t, err)

	assert.Equal(t, resp.AWBNumber, capturedEvent.AWBNumber,
		"Kafka event AWB must match the response AWB")
	assert.Equal(t, resp.TransactionID, capturedEvent.TransactionID)
	assert.Equal(t, "Bandung", capturedEvent.DestCity)
}

// ─────────────────────────────────────────────────────────────────────────────
// GetOrderByAWB Tests
// ─────────────────────────────────────────────────────────────────────────────

// TestGetOrderByAWB_Success verifies that a valid AWB returns the correct order.
func TestGetOrderByAWB_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)
	ctx := context.Background()

	expectedOrder := &model.Order{
		AWBNumber:    "JNE-abc12345",
		Status:       model.StatusOrderCreated,
		SenderName:   "Budi Santoso",
		ReceiverName: "Ani Rahayu",
		TotalPrice:   17000,
	}

	mockRepo.EXPECT().
		FindByAWB(ctx, "JNE-abc12345").
		Return(expectedOrder, nil)

	result, err := svc.GetOrderByAWB(ctx, "JNE-abc12345")

	require.NoError(t, err)
	assert.Equal(t, "JNE-abc12345", result.AWBNumber)
	assert.Equal(t, model.StatusOrderCreated, result.Status)
}

// TestGetOrderByAWB_NotFound verifies that a missing AWB returns an error.
func TestGetOrderByAWB_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPricing := mock.NewMockPricingClient(ctrl)
	mockKafka := mock.NewMockKafkaProducer(ctrl)

	svc := service.NewOrderService(mockRepo, mockPricing, mockKafka)
	ctx := context.Background()

	mockRepo.EXPECT().
		FindByAWB(ctx, "JNE-notexist").
		Return(nil, errors.New("record not found"))

	result, err := svc.GetOrderByAWB(ctx, "JNE-notexist")

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")
}