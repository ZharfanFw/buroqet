package service_test

import (
	"testing"

	"pricing-service/internal/domain"
	"pricing-service/mocks"
	"pricing-service/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTariff(t *testing.T) {

	mockRepo := mocks.NewMockPricingRepository()

	pricingSvc := service.NewPricingService(mockRepo)

	req := domain.CalculationRequest{
	OriginPostalCode:      "10110",
	DestinationPostalCode: "40115",

	WeightKG: 2,

	LengthCM: 10,
	WidthCM:  10,
	HeightCM: 10,

	ServiceType:  "REG",
	UseInsurance: true,
}

	resp := pricingSvc.CalculateTariff(req)

	assert.Equal(t, 20400.0, resp.Total)
	assert.Equal(t, "2-3 Days", resp.Estimated)
}