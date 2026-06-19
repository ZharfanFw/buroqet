package mocks

type MockPricingRepository struct {
}

func NewMockPricingRepository() *MockPricingRepository {
	return &MockPricingRepository{}
}

func (m *MockPricingRepository) GetBaseRate(
	origin string,
	destination string,
	serviceType string,
) float64 {
	return 10000
}