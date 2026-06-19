package handler
import (
	"encoding/json"
	"net/http"
	"strings"

	"settlement-service/internal/service"
)

type SettlementHandler struct {
	service *service.SettlementService
}

func NewSettlementHandler(service *service.SettlementService) *SettlementHandler {
	return &SettlementHandler{service: service}
}

type CommissionRequest struct {
	CourierID   string  `json:"courier_id"`
	AWB         string  `json:"awb"`
	ServiceType string  `json:"service_type"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *SettlementHandler) ProcessCommission(w http.ResponseWriter, r *http.Request) {
	var req CommissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	err := h.service.ProcessDeliveryCommission(r.Context(), req.CourierID, req.AWB, req.ServiceType)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "commission recorded successfully"})
}

func (h *SettlementHandler) GetCourierEarnings(w http.ResponseWriter, r *http.Request) {
	// Ambil courierID dari URL: /api/v1/couriers/{courierID}/earnings
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid path"})
		return
	}
	courierID := parts[4]

	summary, err := h.service.GetCourierEarnings(r.Context(), courierID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

func (h *SettlementHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}