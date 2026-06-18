// warehouse-service/internal/handler/warehouse_handler.go
// HTTP Handler layer untuk Warehouse Management Service.
// Bertanggung jawab: decode request, validasi method HTTP, panggil service, encode response.
// Mengikuti pola yang sama dengan tracking-service/handler (concrete service dependency).

package handler

import (
	"encoding/json"
	"net/http"

	"warehouse-service/internal/service"
)

// WarehouseHandler menangani semua HTTP request untuk WMS
type WarehouseHandler struct {
	service *service.WarehouseService
}

// NewWarehouseHandler membuat instance handler baru
func NewWarehouseHandler(svc *service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: svc}
}

// InboundRequest adalah request body untuk POST /api/v1/inbound
type InboundRequest struct {
	AWB   string `json:"awb"`
	HubID string `json:"hub_id"`
}

// DispatchRequest adalah request body untuk POST /api/v1/dispatch
type DispatchRequest struct {
	ManifestID string `json:"manifest_id"`
}

// writeJSON adalah helper untuk menulis response JSON secara konsisten
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError adalah helper untuk response error yang konsisten
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// =========================================================
// POST /api/v1/inbound
// Proses paket yang masuk ke gudang.
// Body: { "awb": "BQ-2024-JKT-001", "hub_id": "HUB-JKT-01" }
// Response 201: Package yang berhasil di-scan masuk
// Side effect: publish event package.arrived ke Kafka
// =========================================================

func (h *WarehouseHandler) ProcessInbound(w http.ResponseWriter, r *http.Request) {
	// Hanya terima method POST
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req InboundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	// Validasi field wajib di handler layer (validasi bisnis ada di service)
	if req.AWB == "" {
		writeError(w, http.StatusBadRequest, "awb wajib diisi")
		return
	}
	if req.HubID == "" {
		writeError(w, http.StatusBadRequest, "hub_id wajib diisi")
		return
	}

	pkg, err := h.service.ProcessInbound(r.Context(), req.AWB, req.HubID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, pkg)
}

// =========================================================
// POST /api/v1/dispatch
// Dispatch manifest (kumpulan paket) untuk dikirim ke hub tujuan.
// Body: { "manifest_id": "MNF-..." }
// Response 200: pesan sukses
// Side effect: publish event manifest.dispatched ke Kafka, update semua AWB ke ON_TRANSIT
// =========================================================

func (h *WarehouseHandler) DispatchManifest(w http.ResponseWriter, r *http.Request) {
	// Hanya terima method POST
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req DispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if req.ManifestID == "" {
		writeError(w, http.StatusBadRequest, "manifest_id wajib diisi")
		return
	}

	if err := h.service.DispatchManifest(r.Context(), req.ManifestID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "manifest dispatched successfully"})
}

// =========================================================
// GET /health — Liveness probe
// GET /ready  — Readiness probe
// =========================================================

func (h *WarehouseHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "warehouse-service",
	})
}