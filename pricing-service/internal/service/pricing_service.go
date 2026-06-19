package service

import "pricing-service/internal/domain"

type PricingRepository interface {
	GetBaseRate(origin, destination, serviceType string) float64
}

type PricingService struct {
	repo PricingRepository
}

func NewPricingService(repo PricingRepository) *PricingService {
	return &PricingService{
		repo: repo,
	}
}

func (s *PricingService) CalculateTariff(
	req domain.CalculationRequest,
) domain.CalculationResponse {

	actualWeight := req.WeightKG

	volumetric := (req.LengthCM * req.WidthCM * req.HeightCM) / 6000

	finalWeight := actualWeight

	if volumetric > actualWeight {
		finalWeight = volumetric
	}

	baseRate := s.repo.GetBaseRate(
		req.OriginPostalCode,
		req.DestinationPostalCode,
		req.ServiceType,
	)

	base := finalWeight * baseRate
	insurance := 0.0

	if req.UseInsurance {
		insurance = base * 0.02
	}

	total := base + insurance

	return domain.CalculationResponse{
		BaseTariff: base,
		Insurance:  insurance,
		Discount:   0,
		Total:      total,
		Estimated:  "2-3 Days",
	}
}