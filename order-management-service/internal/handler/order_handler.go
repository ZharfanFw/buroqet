package handler

import (
	"net/http"

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
    // Health check untuk Kubernetes probe
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    v1 := r.Group("/api/v1")
    {
        v1.POST("/orders", h.CreateOrder)
        v1.GET("/orders/:awb", h.GetOrderByAWB)
    }
}

// CreateOrder godoc
// POST /api/v1/orders
// Accepts a CreateOrderRequest JSON body, creates the order, and returns a 201 with the order details.
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req model.CreateOrderRequest
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