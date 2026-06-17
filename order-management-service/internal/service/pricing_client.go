package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func (c *httpPricingClient) GetPrice(ctx context.Context, req model.PricingRequest) (*model.PricingResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/pricing", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("pricing service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pricing service returned status %d", resp.StatusCode)
	}

	var pricingResp model.PricingResponse
	if err := json.NewDecoder(resp.Body).Decode(&pricingResp); err != nil {
		return nil, err
	}
	return &pricingResp, nil
}

// stubPricingClient is used when no real Pricing Service URL is configured.
// It returns a hardcoded price so the rest of the OMS can be developed/tested independently.
type stubPricingClient struct{}

func (s *stubPricingClient) GetPrice(_ context.Context, req model.PricingRequest) (*model.PricingResponse, error) {
	baseFare := 15000.0
	if req.ServiceType == model.ServiceExpress {
		baseFare = 30000.0
	}
	return &model.PricingResponse{
		BaseFare:     baseFare,
		Insurance:    2000.0,
		Discount:     0,
		TotalPrice:   baseFare + 2000.0,
		EstimatedSLA: "2-3 hari",
	}, nil
}