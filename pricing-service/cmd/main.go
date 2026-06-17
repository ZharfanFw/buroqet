package main

import (
	"pricing-service/internal/handler"
	"pricing-service/mocks"
	"pricing-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	repo := mocks.NewMockPricingRepository()

	pricingService := service.NewPricingService(repo)

	pricingHandler := handler.NewPricingHandler(pricingService)

	router.POST("/pricing/calculate", pricingHandler.CalculatePricing)

	router.Run(":8081")
}