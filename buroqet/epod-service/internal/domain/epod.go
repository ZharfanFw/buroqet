package domain

type UploadRequest struct {
	AWB       string
	CourierID string
	Latitude  float64
	Longitude float64
	FileName  string
}

type UploadResponse struct {
	Status   string
	ImageURL string
}