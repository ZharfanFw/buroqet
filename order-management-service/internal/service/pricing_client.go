package service

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"order-management-service/internal/domain"
)

// httpPricingClient calls the real Pricing & Routing Service over HTTP.
type httpPricingClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPPricingClient returns a PricingClient backed by HTTP.
// If baseURL is empty, it falls back to the stub implementation so the service
// can still start without a running Pricing Service (useful during local dev).
func NewHTTPPricingClient(baseURL string) PricingClient {
	if baseURL == "" {
		return &stubPricingClient{}
	}
	return &httpPricingClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type externalCalculationRequest struct {
	OriginPostalCode      string  `json:"origin_postal_code"`
	DestinationPostalCode string  `json:"destination_postal_code"`
	WeightKG              float64 `json:"weight_kg"`
	LengthCM              float64 `json:"length_cm"`
	WidthCM               float64 `json:"width_cm"`
	HeightCM              float64 `json:"height_cm"`
	ServiceType           string  `json:"service_type"`
	UseInsurance          bool    `json:"use_insurance"`
}

type externalCalculationResponse struct {
	BaseTariff float64 `json:"base_tariff"`
	Insurance  float64 `json:"insurance"`
	Discount   float64 `json:"discount"`
	Total      float64 `json:"total"`
	Estimated  string  `json:"estimated"`
}

type externalPricingResponseWrapper struct {
	Message string                      `json:"message"`
	Data    externalCalculationResponse `json:"data"`
}

func (c *httpPricingClient) GetPrice(ctx context.Context, req domain.PricingRequest) (*domain.PricingResponse, error) {
	extReq := externalCalculationRequest{
		OriginPostalCode:      req.OriginPostal,
		DestinationPostalCode: req.DestPostal,
		WeightKG:              req.Weight,
		LengthCM:              req.Length,
		WidthCM:               req.Width,
		HeightCM:              req.Height,
		ServiceType:           string(req.ServiceType),
		UseInsurance:          false,
	}

	body, err := json.Marshal(extReq)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/pricing/calculate", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("WARNING: pricing service unreachable (%v). Falling back to stub pricing.", err)
		stub := &stubPricingClient{}
		return stub.GetPrice(ctx, req)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("WARNING: pricing service returned status %d. Falling back to stub pricing.", resp.StatusCode)
		stub := &stubPricingClient{}
		return stub.GetPrice(ctx, req)
	}

	var wrapper externalPricingResponseWrapper
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}
	return &domain.PricingResponse{
		BaseFare:     wrapper.Data.BaseTariff,
		Insurance:    wrapper.Data.Insurance,
		Discount:     wrapper.Data.Discount,
		TotalPrice:   wrapper.Data.Total,
		EstimatedSLA: wrapper.Data.Estimated,
	}, nil
}

// stubPricingClient is used when no real Pricing Service URL is configured.
// It returns a hardcoded price so the rest of the OMS can be developed/tested independently.
type stubPricingClient struct{}

func (s *stubPricingClient) GetPrice(_ context.Context, req domain.PricingRequest) (*domain.PricingResponse, error) {
	baseFare := 15000.0
	if req.ServiceType == domain.ServiceExpress {
		baseFare = 30000.0
	}
	return &domain.PricingResponse{
		BaseFare:     baseFare,
		Insurance:    2000.0,
		Discount:     0,
		TotalPrice:   baseFare + 2000.0,
		EstimatedSLA: "2-3 hari",
	}, nil
}
