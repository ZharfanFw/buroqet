package functional_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"order-management-service/internal/domain"
	"order-management-service/internal/handler"
	"order-management-service/internal/kafka"
	"order-management-service/internal/repository"
	"order-management-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────────────────────────────────────
// Test Suite Setup
// ─────────────────────────────────────────────────────────────────────────────

// testApp holds all wired-up components for the functional test suite.
// The DB is a real PostgreSQL connection; Kafka and Pricing are stubs
// so the test environment does not need a running broker or pricing service.
type testApp struct {
	router *gin.Engine
	db     *gorm.DB
}

// setupTestApp initialises a real DB connection and wires up the full
// handler → service → repository stack.
//
// Environment variables (set these in your CI/Jenkins environment or .env):
//
//	TEST_DATABASE_URL  postgresql DSN for the test DB
//	                   default: "host=localhost user=postgres password=postgres dbname=oms_test_db port=5432 sslmode=disable"
func setupTestApp(t *testing.T) *testApp {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=oms_test_db port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Gagal konek ke database test: %v", err)
	}

	// Auto-migrate creates/updates the orders table for the test database
	require.NoError(t, db.AutoMigrate(&domain.Order{}))

	// Use stub implementations so functional tests only need a DB,
	// not a running Kafka broker or Pricing Service
	orderRepo := repository.NewOrderRepository(db)
	pricingClient := newStubPricingClient()
	kafkaProducer := newStubKafkaProducer()

	orderSvc := service.NewOrderService(orderRepo, pricingClient, kafkaProducer)
	orderHandler := handler.NewOrderHandler(orderSvc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	orderHandler.RegisterRoutes(r)

	return &testApp{router: r, db: db}
}

// cleanupOrders removes all rows from the orders table after each test
// so tests are isolated from each other.
func (a *testApp) cleanupOrders(t *testing.T) {
	t.Helper()
	a.db.Exec("DELETE FROM orders")
}

// ─────────────────────────────────────────────────────────────────────────────
// Stub implementations (replace external dependencies in functional tests)
// ─────────────────────────────────────────────────────────────────────────────

// stubPricingClient returns a fixed price so functional tests do not need
// the real Pricing & Routing Service to be running.
type stubPricingClient struct{}

func newStubPricingClient() service.PricingClient {
	return &stubPricingClient{}
}

func (s *stubPricingClient) GetPrice(_ context.Context, req domain.PricingRequest) (*domain.PricingResponse, error) {
	baseFare := 15000.0
	if req.ServiceType == domain.ServiceExpress {
		baseFare = 30000.0
	}
	return &domain.PricingResponse{
		BaseFare:     baseFare,
		Insurance:    2000,
		Discount:     0,
		TotalPrice:   baseFare + 2000,
		EstimatedSLA: "2-3 hari",
	}, nil
}

// stubKafkaProducer silently swallows all Kafka publish calls so functional
// tests do not need a running Kafka broker.
type stubKafkaProducer struct{}

func newStubKafkaProducer() kafka.Producer {
	return &stubKafkaProducer{}
}

func (s *stubKafkaProducer) PublishOrderCreated(_ context.Context, _ kafka.OrderCreatedEvent) error {
	return nil
}

func (s *stubKafkaProducer) Close() error {
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func validOrderPayload() map[string]interface{} {
	return map[string]interface{}{
		"sender_name":      "Budi Santoso",
		"sender_phone":     "081234567890",
		"sender_address":   "Jl. Merdeka No.1, Jakarta",
		"origin_city":      "Jakarta",
		"origin_postal":    "10110",
		"receiver_name":    "Ani Rahayu",
		"receiver_phone":   "089876543210",
		"receiver_address": "Jl. Sudirman No.5, Bandung",
		"dest_city":        "Bandung",
		"dest_postal":      "40111",
		"weight_actual":    2.0,
		"length":           20.0,
		"width":            15.0,
		"height":           10.0,
		"service_type":     "REGULER",
		"payment_type":     "NON_COD",
	}
}

func toJSON(t *testing.T, v interface{}) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

func doRequest(router *gin.Engine, method, path string, body *bytes.Buffer) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, path, body)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	router.ServeHTTP(w, req)
	return w
}

// ─────────────────────────────────────────────────────────────────────────────
// Functional Tests
// ─────────────────────────────────────────────────────────────────────────────

// FT-01: Creating a valid order must persist it to the database and return 201.
func TestFunctional_CreateOrder_PersistedToDatabase(t *testing.T) {
	app := setupTestApp(t)
	defer app.cleanupOrders(t)

	w := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, validOrderPayload()))

	// Assert HTTP response
	assert.Equal(t, http.StatusCreated, w.Code)

	var respBody map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &respBody))
	assert.True(t, respBody["success"].(bool))

	data := respBody["data"].(map[string]interface{})
	awb := data["awb_number"].(string)
	assert.NotEmpty(t, awb)

	// Assert the order is actually in the database
	var order domain.Order
	err := app.db.Where("awb_number = ?", awb).First(&order).Error
	require.NoError(t, err, "Order must be persisted to the database")

	assert.Equal(t, awb, order.AWBNumber)
	assert.Equal(t, domain.StatusOrderCreated, order.Status)
	assert.Equal(t, "Budi Santoso", order.SenderName)
	assert.Equal(t, "Ani Rahayu", order.ReceiverName)
	assert.Equal(t, domain.ServiceRegular, order.ServiceType)
	assert.Equal(t, domain.PaymentNonCOD, order.PaymentType)
	assert.Greater(t, order.TotalPrice, 0.0)
}

// FT-02: Retrieving an order by AWB must return the persisted data correctly.
func TestFunctional_GetOrderByAWB_ReturnsPersistedData(t *testing.T) {
	app := setupTestApp(t)
	defer app.cleanupOrders(t)

	// First, create an order
	w := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, validOrderPayload()))
	require.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &createResp))
	awb := createResp["data"].(map[string]interface{})["awb_number"].(string)

	// Then, retrieve it by AWB
	w2 := doRequest(app.router, http.MethodGet, "/api/v1/orders/"+awb, nil)

	assert.Equal(t, http.StatusOK, w2.Code)

	var getResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &getResp))
	assert.True(t, getResp["success"].(bool))

	data := getResp["data"].(map[string]interface{})
	assert.Equal(t, awb, data["awb_number"])
	assert.Equal(t, string(domain.StatusOrderCreated), data["status"])
	assert.Equal(t, "Budi Santoso", data["sender_name"])
}

// FT-03: Two orders must receive different AWB numbers (uniqueness constraint).
func TestFunctional_CreateOrder_AWBIsUniquePerOrder(t *testing.T) {
	app := setupTestApp(t)
	defer app.cleanupOrders(t)

	w1 := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, validOrderPayload()))
	w2 := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, validOrderPayload()))

	require.Equal(t, http.StatusCreated, w1.Code)
	require.Equal(t, http.StatusCreated, w2.Code)

	var r1, r2 map[string]interface{}
	require.NoError(t, json.Unmarshal(w1.Body.Bytes(), &r1))
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &r2))

	awb1 := r1["data"].(map[string]interface{})["awb_number"].(string)
	awb2 := r2["data"].(map[string]interface{})["awb_number"].(string)

	assert.NotEqual(t, awb1, awb2, "Each order must get a unique AWB number")
}

// FT-04: Querying a non-existent AWB must return 404.
func TestFunctional_GetOrderByAWB_NotFound(t *testing.T) {
	app := setupTestApp(t)
	defer app.cleanupOrders(t)

	w := doRequest(app.router, http.MethodGet, "/api/v1/orders/JNE-notexist99", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.False(t, body["success"].(bool))
}

// FT-05: Express service order should have higher price than Reguler.
func TestFunctional_CreateOrder_ExpressIsMoreExpensiveThanReguler(t *testing.T) {
	app := setupTestApp(t)
	defer app.cleanupOrders(t)

	// Create REGULER order
	regularPayload := validOrderPayload()
	wReg := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, regularPayload))
	require.Equal(t, http.StatusCreated, wReg.Code)

	// Create EXPRESS order
	expressPayload := validOrderPayload()
	expressPayload["service_type"] = "EXPRESS"
	wExp := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, expressPayload))
	require.Equal(t, http.StatusCreated, wExp.Code)

	var regResp, expResp map[string]interface{}
	require.NoError(t, json.Unmarshal(wReg.Body.Bytes(), &regResp))
	require.NoError(t, json.Unmarshal(wExp.Body.Bytes(), &expResp))

	regPrice := regResp["data"].(map[string]interface{})["total_price"].(float64)
	expPrice := expResp["data"].(map[string]interface{})["total_price"].(float64)

	assert.Greater(t, expPrice, regPrice, "EXPRESS service must cost more than REGULER")
}

// FT-06: COD order must NOT have a payment URL.
func TestFunctional_CreateOrder_COD_NoPaymentURL(t *testing.T) {
	app := setupTestApp(t)
	defer app.cleanupOrders(t)

	payload := validOrderPayload()
	payload["payment_type"] = "COD"

	w := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, payload))
	require.Equal(t, http.StatusCreated, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	data := body["data"].(map[string]interface{})
	paymentURL, exists := data["payment_url"]
	assert.True(t, !exists || paymentURL == "" || paymentURL == nil,
		"COD orders must not have a payment URL")
}

// FT-07: NON-COD order must have a payment URL.
func TestFunctional_CreateOrder_NonCOD_HasPaymentURL(t *testing.T) {
	app := setupTestApp(t)
	defer app.cleanupOrders(t)

	w := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, validOrderPayload()))
	require.Equal(t, http.StatusCreated, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	data := body["data"].(map[string]interface{})
	assert.NotEmpty(t, data["payment_url"], "NON-COD orders must have a payment URL")
}

// FT-08: Volumetric weight must be correctly stored in the database.
func TestFunctional_CreateOrder_VolumetricWeightStoredCorrectly(t *testing.T) {
	app := setupTestApp(t)
	defer app.cleanupOrders(t)

	payload := validOrderPayload()
	// 30 x 20 x 10 / 6000 = 1.0 kg volumetric
	payload["length"] = 30.0
	payload["width"] = 20.0
	payload["height"] = 10.0

	w := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, payload))
	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	awb := resp["data"].(map[string]interface{})["awb_number"].(string)

	var order domain.Order
	require.NoError(t, app.db.Where("awb_number = ?", awb).First(&order).Error)

	expectedVolumetric := (30.0 * 20.0 * 10.0) / 6000.0
	assert.InDelta(t, expectedVolumetric, order.WeightVolumetri, 0.001,
		"Volumetric weight (L*W*H/6000) must be stored correctly in the DB")
}

// FT-09: Invalid payload must not create any database record.
func TestFunctional_CreateOrder_InvalidPayload_NoDBRecord(t *testing.T) {
	app := setupTestApp(t)
	defer app.cleanupOrders(t)

	badPayload := map[string]interface{}{
		"sender_name": "Hanya Nama Saja",
		// semua field lain tidak ada
	}

	w := doRequest(app.router, http.MethodPost, "/api/v1/orders", toJSON(t, badPayload))
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify no record was inserted
	var count int64
	app.db.Model(&domain.Order{}).Count(&count)
	assert.Equal(t, int64(0), count, "No order must be saved when the request payload is invalid")
}
