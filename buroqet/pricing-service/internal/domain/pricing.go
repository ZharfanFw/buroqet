package domain

type CalculationRequest struct {
	OriginPostalCode      string  `json:"origin_postal_code"`
	DestinationPostalCode string  `json:"destination_postal_code"`

	OriginLat float64 `json:"origin_lat"`
	OriginLon float64 `json:"origin_lon"`

	DestinationLat float64 `json:"destination_lat"`
	DestinationLon float64 `json:"destination_lon"`

	WeightKG float64 `json:"weight_kg"`

	LengthCM float64 `json:"length_cm"`
	WidthCM  float64 `json:"width_cm"`
	HeightCM float64 `json:"height_cm"`

	ServiceType  string `json:"service_type"`
	UseInsurance bool   `json:"use_insurance"`
	PromoCode    string `json:"promo_code"`
}

type CalculationResponse struct {
	BaseTariff float64 `json:"base_tariff"`
	Insurance  float64 `json:"insurance"`
	Discount   float64 `json:"discount"`
	Total      float64 `json:"total"`
	Estimated  string  `json:"estimated"`
}