package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-management-service/internal/handler"
	"order-management-service/internal/domain"
	"order-management-service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// setupRouter creates a fresh Gin router in test mode with the handler registered.
// Menggunakan gin.TestMode agar output gin tidak berisik saat test berjalan.
func setupRouter(h *handler.OrderHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h.RegisterRoutes(r)
	return r
}

// validRequestBody returns a JSON body that passes all binding validations.
func validRequestBody() map[string]interface{} {
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

// toJSONReader marshals a map to a JSON reader suitable for http.NewRequest.
func toJSONReader(t *testing.T, body interface{}) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /api/v1/orders
// ─────────────────────────────────────────────────────────────────────────────

// TestCreateOrderHandler_Success verifies that a valid request returns 201
// with the correct AWB, transaction ID, and status in the response body.
func TestCreateOrderHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockOrderService(ctrl)
	h := handler.NewOrderHandler(mockSvc)
	router := setupRouter(h)

	expectedResp := &model.CreateOrderResponse{
		AWBNumber:     "JNE-abc12345",
		TransactionID: "txn-uuid-001",
		Status:        model.StatusOrderCreated,
		TotalPrice:    17000,
		PaymentURL:    "https://pay.example.com/invoice/txn-uuid-001",
	}

	mockSvc.EXPECT().
		CreateOrder(gomock.Any(), gomock.Any()).
		Return(expectedResp, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", toJSONReader(t, validRequestBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	assert.True(t, body["success"].(bool))
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "JNE-abc12345", data["awb_number"])
	assert.Equal(t, "txn-uuid-001", data["transaction_id"])
	assert.Equal(t, string(model.StatusOrderCreated), data["status"])
	assert.NotEmpty(t, data["payment_url"])
}

// TestCreateOrderHandler_InvalidBody verifies that a malformed / incomplete JSON body
// returns 400 Bad Request WITHOUT calling the service layer.
func TestCreateOrderHandler_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockOrderService(ctrl)
	// Service must not be called at all when binding fails
	mockSvc.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(0)

	h := handler.NewOrderHandler(mockSvc)
	router := setupRouter(h)

	incompleteBody := map[string]interface{}{
		"sender_name": "Budi Santoso",
		// Missing all required fields
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", toJSONReader(t, incompleteBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.False(t, body["success"].(bool))
	assert.NotEmpty(t, body["error"])
}

// TestCreateOrderHandler_InvalidServiceType verifies that an invalid enum value
// for service_type (not REGULER or EXPRESS) is rejected at the binding layer.
func TestCreateOrderHandler_InvalidServiceType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockOrderService(ctrl)
	mockSvc.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(0)

	h := handler.NewOrderHandler(mockSvc)
	router := setupRouter(h)

	body := validRequestBody()
	body["service_type"] = "SUPER_KILAT" // invalid enum

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", toJSONReader(t, body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestCreateOrderHandler_ZeroWeight verifies that weight_actual = 0 fails validation
// (the binding tag requires gt=0).
func TestCreateOrderHandler_ZeroWeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockOrderService(ctrl)
	mockSvc.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(0)

	h := handler.NewOrderHandler(mockSvc)
	router := setupRouter(h)

	body := validRequestBody()
	body["weight_actual"] = 0 // invalid: must be > 0

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", toJSONReader(t, body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestCreateOrderHandler_ServiceError verifies that when the service layer returns
// an error, the handler responds with 500 Internal Server Error.
func TestCreateOrderHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockOrderService(ctrl)
	mockSvc.EXPECT().
		CreateOrder(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("pricing service timeout"))

	h := handler.NewOrderHandler(mockSvc)
	router := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", toJSONReader(t, validRequestBody()))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.False(t, body["success"].(bool))
	assert.Contains(t, body["error"].(string), "pricing service timeout")
}

// TestCreateOrderHandler_EmptyBody verifies that sending an empty body returns 400.
func TestCreateOrderHandler_EmptyBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockOrderService(ctrl)
	mockSvc.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Times(0)

	h := handler.NewOrderHandler(mockSvc)
	router := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /api/v1/orders/:awb
// ─────────────────────────────────────────────────────────────────────────────

// TestGetOrderByAWBHandler_Success verifies that a valid AWB returns 200
// with the full order details.
func TestGetOrderByAWBHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockOrderService(ctrl)
	h := handler.NewOrderHandler(mockSvc)
	router := setupRouter(h)

	expectedOrder := &model.Order{
		AWBNumber:    "JNE-abc12345",
		Status:       model.StatusOrderCreated,
		SenderName:   "Budi Santoso",
		ReceiverName: "Ani Rahayu",
		TotalPrice:   17000,
	}

	mockSvc.EXPECT().
		GetOrderByAWB(gomock.Any(), "JNE-abc12345").
		Return(expectedOrder, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders/JNE-abc12345", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.True(t, body["success"].(bool))

	data := body["data"].(map[string]interface{})
	assert.Equal(t, "JNE-abc12345", data["awb_number"])
	assert.Equal(t, string(model.StatusOrderCreated), data["status"])
}

// TestGetOrderByAWBHandler_NotFound verifies that an unknown AWB returns 404.
func TestGetOrderByAWBHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockOrderService(ctrl)
	h := handler.NewOrderHandler(mockSvc)
	router := setupRouter(h)

	mockSvc.EXPECT().
		GetOrderByAWB(gomock.Any(), "JNE-notexist").
		Return(nil, errors.New("order not found: record not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders/JNE-notexist", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.False(t, body["success"].(bool))
}

// TestGetOrderByAWBHandler_ResponseStructure verifies the exact JSON shape
// returned by a successful GET request.
func TestGetOrderByAWBHandler_ResponseStructure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mock.NewMockOrderService(ctrl)
	h := handler.NewOrderHandler(mockSvc)
	router := setupRouter(h)

	mockSvc.EXPECT().
		GetOrderByAWB(gomock.Any(), gomock.Any()).
		Return(&model.Order{
			AWBNumber:    "JNE-test1234",
			Status:       model.StatusOrderCreated,
			SenderName:   "Sender",
			ReceiverName: "Receiver",
			TotalPrice:   20000,
			ServiceType:  model.ServiceRegular,
			PaymentType:  model.PaymentNonCOD,
		}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders/JNE-test1234", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	// Verify top-level keys exist
	assert.Contains(t, body, "success")
	assert.Contains(t, body, "data")

	// Verify data fields
	data := body["data"].(map[string]interface{})
	assert.Contains(t, data, "awb_number")
	assert.Contains(t, data, "status")
	assert.Contains(t, data, "total_price")
	assert.Contains(t, data, "service_type")
}