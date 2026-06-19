package handler

import (
	"net/http"
	"strconv"

	"order-management-service/internal/domain"
	"order-management-service/internal/service"

	"github.com/gin-gonic/gin"
)

// OrderHandler holds the HTTP handlers for the OMS endpoints.
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler creates a new OrderHandler with the given service.
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// RegisterRoutes mounts the OMS routes onto the given Gin engine.
func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	r.Use(corsMiddleware())

	// Health check untuk Kubernetes probe
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 1. Group /api/v1/orders (untuk test dan local dev legacy)
	v1 := r.Group("/api/v1")
	{
		// NOTE: specific sub-paths (/transaction/:id, /customer/:id) must be
		// registered BEFORE the wildcard route (/orders/:awb) to avoid Gin conflicts.
		v1.GET("/orders", h.ListOrders)
		v1.POST("/orders", h.CreateOrder)
		v1.GET("/orders/transaction/:id", h.GetOrderByTransactionID)
		v1.GET("/orders/customer/:id", h.GetOrdersByCustomerID)
		v1.GET("/orders/:awb", h.GetOrderByAWB)
		v1.PATCH("/orders/:awb/status", h.UpdateOrderStatus)
		v1.DELETE("/orders/:awb", h.CancelOrder)
	}

	// 2. Group root / (untuk Ingress rewrite target di AKS)
	root := r.Group("")
	{
		// NOTE: specific sub-paths (/transaction/:id, /customer/:id) must be
		// registered BEFORE the wildcard route (/:awb) to avoid Gin conflicts.
		root.GET("/", h.ListOrders)
		root.POST("/", h.CreateOrder)
		root.GET("/transaction/:id", h.GetOrderByTransactionID)
		root.GET("/customer/:id", h.GetOrdersByCustomerID)
		root.GET("/:awb", h.GetOrderByAWB)
		root.PATCH("/:awb/status", h.UpdateOrderStatus)
		root.DELETE("/:awb", h.CancelOrder)
	}
}

// CreateOrder godoc
// POST /api/v1/orders
// Accepts a CreateOrderRequest JSON body, creates the order, and returns 201.
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	resp, err := h.orderService.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    resp,
	})
}

// GetOrderByAWB godoc
// GET /api/v1/orders/:awb
// Returns the full order details for the given AWB number.
func (h *OrderHandler) GetOrderByAWB(c *gin.Context) {
	awb := c.Param("awb")
	if awb == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "awb parameter is required",
		})
		return
	}

	order, err := h.orderService.GetOrderByAWB(c.Request.Context(), awb)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

// ListOrders godoc
// GET /api/v1/orders?status=ORDER_CREATED&page=1&limit=10
// Returns a paginated list of orders, optionally filtered by status or customer_id.
func (h *OrderHandler) ListOrders(c *gin.Context) {
	var req domain.ListOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	result, err := h.orderService.ListOrders(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateOrderStatus godoc
// PATCH /api/v1/orders/:awb/status
// Updates the lifecycle status of an existing order.
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	awb := c.Param("awb")
	if awb == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "awb parameter is required",
		})
		return
	}

	var req domain.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	order, err := h.orderService.UpdateOrderStatus(c.Request.Context(), awb, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

// GetOrderByTransactionID godoc
// GET /api/v1/orders/transaction/:id
// Finds an order by its payment reference / transaction ID.
func (h *OrderHandler) GetOrderByTransactionID(c *gin.Context) {
	transactionID := c.Param("id")
	if transactionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "transaction id parameter is required",
		})
		return
	}

	order, err := h.orderService.GetOrderByTransactionID(c.Request.Context(), transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

// GetOrdersByCustomerID godoc
// GET /api/v1/orders/customer/:id?page=1&limit=10
// Returns a paginated list of orders for a specific customer.
func (h *OrderHandler) GetOrdersByCustomerID(c *gin.Context) {
	customerID := c.Param("id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "customer id parameter is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.orderService.GetOrdersByCustomerID(c.Request.Context(), customerID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// CancelOrder godoc
// DELETE /api/v1/orders/:awb
// Cancels (deletes) an order only if its status is still ORDER_CREATED.
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	awb := c.Param("awb")
	if awb == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "awb parameter is required",
		})
		return
	}

	resp, err := h.orderService.CancelOrder(c.Request.Context(), awb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    resp,
	})
}

// corsMiddleware adds CORS headers to Gin requests
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
