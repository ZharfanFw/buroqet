// settlement-service/internal/handler/settlement_handler.go
// HTTP Handler layer untuk Settlement & Commission Service.
// Bertanggung jawab: decode request, validasi method HTTP, panggil service, encode response.

package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"settlement-service/internal/service"
)

// SettlementHandler menangani semua HTTP request untuk Settlement Service
type SettlementHandler struct {
	service *service.SettlementService
}

// NewSettlementHandler membuat instance handler baru
func NewSettlementHandler(svc *service.SettlementService) *SettlementHandler {
	return &SettlementHandler{service: svc}
}

// CommissionRequest adalah request body untuk POST /api/v1/commissions
type CommissionRequest struct {
	CourierID   string `json:"courier_id"`
	AWB         string `json:"awb"`
	ServiceType string `json:"service_type"`
}

// writeJSON adalah helper untuk menulis response JSON secara konsisten
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError adalah helper untuk response error yang konsisten
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// =========================================================
// CommissionsHandler — dispatches GET and POST for /api/v1/commissions
// GET  → list all commissions
// POST → record new commission
// =========================================================

func (h *SettlementHandler) CommissionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.ListCommissions(w, r)
		return
	}
	h.ProcessCommission(w, r)
}

// =========================================================
// POST /api/v1/commissions
// Proses dan catat komisi kurir secara manual.
// Body: { "courier_id": "...", "awb": "...", "service_type": "REGULER|EXPRESS" }
// Response 201: pesan sukses
// Digunakan juga sebagai fallback jika Kafka consumer tidak bisa auto-proses.
// =========================================================

func (h *SettlementHandler) ProcessCommission(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Hanya terima method POST
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CommissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	// Validasi field wajib di handler layer
	if req.CourierID == "" {
		writeError(w, http.StatusBadRequest, "courier_id wajib diisi")
		return
	}
	if req.AWB == "" {
		writeError(w, http.StatusBadRequest, "awb wajib diisi")
		return
	}
	if req.ServiceType == "" {
		req.ServiceType = "REGULER" // default jika tidak diisi
	}

	err := h.service.ProcessDeliveryCommission(r.Context(), req.CourierID, req.AWB, req.ServiceType)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "commission recorded successfully"})
}

// =========================================================
// GET /api/v1/commissions
// Ambil daftar semua commission log.
// Response 200: array of CommissionLog
// =========================================================

func (h *SettlementHandler) ListCommissions(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	logs, err := h.service.GetAllCommissions(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "gagal mengambil data commissions: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, logs)
}

// =========================================================
// GET /api/v1/couriers/{courierID}/earnings
// Ambil rekap penghasilan kurir tertentu.
// Response 200: CourierSummary (total, pending, paid)
// URL: /api/v1/couriers/COURIER-001/earnings
// =========================================================

func (h *SettlementHandler) GetCourierEarnings(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Hanya terima method GET
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Ekstrak courierID dari URL: /api/v1/couriers/{courierID}/earnings
	// Path format: ["", "api", "v1", "couriers", "{courierID}", "earnings"]
	courierID := extractCourierIDFromPath(r.URL.Path)
	if courierID == "" {
		writeError(w, http.StatusBadRequest, "courier ID tidak ditemukan di URL")
		return
	}

	summary, err := h.service.GetCourierEarnings(r.Context(), courierID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

// =========================================================
// GET /health — Liveness probe
// GET /ready  — Readiness probe
// =========================================================

func (h *SettlementHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "settlement-service",
	})
}

// =========================================================
// HELPER
// =========================================================

// extractCourierIDFromPath mengekstrak courierID dari URL path.
// Contoh: "/api/v1/couriers/COURIER-001/earnings" → "COURIER-001"
func extractCourierIDFromPath(path string) string {
	// Hapus trailing slash jika ada
	path = strings.TrimRight(path, "/")
	parts := strings.Split(path, "/")
	// Format: ["", "api", "v1", "couriers", "{courierID}", "earnings"]
	// Index:    0     1     2      3              4              5
	if len(parts) >= 5 {
		return parts[4]
	}
	return ""
}