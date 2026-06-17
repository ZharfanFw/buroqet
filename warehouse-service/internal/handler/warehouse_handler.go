package handler
import (
	"encoding/json"
	"net/http"

	"warehouse-service/internal/service"
)

type WarehouseHandler struct {
	service *service.WarehouseService
}

func NewWarehouseHandler(service *service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

type InboundRequest struct {
	AWB   string `json:"awb"`
	HubID string `json:"hub_id"`
}

type DispatchRequest struct {
	ManifestID string `json:"manifest_id"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *WarehouseHandler) ProcessInbound(w http.ResponseWriter, r *http.Request) {
	var req InboundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	pkg, err := h.service.ProcessInbound(r.Context(), req.AWB, req.HubID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, pkg)
}

func (h *WarehouseHandler) DispatchManifest(w http.ResponseWriter, r *http.Request) {
	var req DispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.DispatchManifest(r.Context(), req.ManifestID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "manifest dispatched successfully"})
}

func (h *WarehouseHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}