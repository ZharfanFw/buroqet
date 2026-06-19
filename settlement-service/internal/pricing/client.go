package pricing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type pricingClient struct {
	baseURL string
	client  *http.Client
}

func NewPricingClient(baseURL string) *pricingClient {
	return &pricingClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type rateResponse struct {
	Rate float64 `json:"rate"`
}

func (p *pricingClient) GetCommissionRate(ctx context.Context, serviceType string) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/commission-rate?type=%s", p.baseURL, serviceType)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		// Fallback rate kalau pricing service tidak bisa diakses
		// Ini common pattern di microservices
		return p.fallbackRate(serviceType), nil
	}
	defer resp.Body.Close()

	var result rateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return p.fallbackRate(serviceType), nil
	}

	return result.Rate, nil
}

func (p *pricingClient) fallbackRate(serviceType string) float64 {
	rates := map[string]float64{
		"REGULER": 3500,
		"EXPRESS": 5000,
		"SAME_DAY": 7500,
	}
	if rate, ok := rates[serviceType]; ok {
		return rate
	}
	return 3500 // default
}