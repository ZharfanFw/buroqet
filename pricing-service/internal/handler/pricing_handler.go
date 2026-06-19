package handler

import (
	"net/http"

	"pricing-service/internal/domain"
	"pricing-service/internal/service"

	"github.com/gin-gonic/gin"
)

type PricingHandler struct {
	service *service.PricingService
}

func NewPricingHandler(service *service.PricingService) *PricingHandler {
	return &PricingHandler{
		service: service,
	}
}

func (h *PricingHandler) CalculatePricing(c *gin.Context) {
	var req domain.CalculationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp := h.service.CalculateTariff(req)

	c.JSON(http.StatusOK, gin.H{
		"message": "pricing calculated successfully",
		"data":    resp,
	})
}