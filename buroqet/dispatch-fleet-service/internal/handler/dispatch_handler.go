package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"dispatch-fleet/internal/domain"
)

// DispatchHandler menangani request HTTP terkait dispatching.
type DispatchHandler struct {
	service domain.DispatchService
}

// NewDispatchHandler membuat instance baru DispatchHandler.
func NewDispatchHandler(service domain.DispatchService) *DispatchHandler {
	return &DispatchHandler{
		service: service,
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