package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"dispatch-fleet/internal/domain"
	"dispatch-fleet/internal/kafka"
)

type DispatchHandler struct {
	service  domain.DispatchService
	producer kafka.Producer
}

// NewDispatchHandler membuat instance baru DispatchHandler.
func NewDispatchHandler(service domain.DispatchService, producer kafka.Producer) *DispatchHandler {
	return &DispatchHandler{
		service:  service,
		producer: producer,
	}
}

// DispatchRequest mendefinisikan format input JSON dari client.
type DispatchRequest struct {
	PickupLat    float64 `json:"pickup_lat"`
	PickupLon    float64 `json:"pickup_lon"`
	RadiusMeters float64 `json:"radius_meters"`
}

// Assign menangani endpoint POST /v1/dispatch/assign.
func (h *DispatchHandler) Assign(w http.ResponseWriter, r *http.Request) {
	// 1. Validasi Method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decode JSON Body
	var req DispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 3. Panggil Business Logic di Service
	// Mengonversi request ke domain.Point
	pickupLoc := domain.Point{
		Longitude: req.PickupLon,
		Latitude:  req.PickupLat,
	}

	result, err := h.service.AssignCourierToPickup(r.Context(), pickupLoc, req.RadiusMeters)
	if err != nil {
		if errors.Is(err, domain.ErrNoCourierAvailable) {
			h.respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// 4. Kirim Respon Sukses
	h.respondWithJSON(w, http.StatusOK, result)
}

// Helper untuk mengirim respon JSON
func (h *DispatchHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Helper untuk mengirim respon Error
func (h *DispatchHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

// StartDeliveryRequest mendefinisikan format input JSON untuk memulai pengiriman.
type StartDeliveryRequest struct {
	AWB      string `json:"awb"`
	Location string `json:"location"`
}

// StartDelivery menangani endpoint POST /v1/dispatch/start-delivery.
func (h *DispatchHandler) StartDelivery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.AWB == "" {
		h.respondWithError(w, http.StatusBadRequest, "awb is required")
		return
	}

	if req.Location == "" {
		req.Location = "Unknown Location"
	}

	err := h.producer.PublishPackageDispatched(r.Context(), req.AWB, req.Location)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Gagal mempublish event package.dispatched")
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Pengiriman dimulai, event package.dispatched berhasil dipublish",
	})
}