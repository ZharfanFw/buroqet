// internal/handler/tracking_handler.go
// HTTP Handler layer untuk Tracking & Status Service.
// Handler menerima request HTTP, decode payload, panggil service, encode response.

package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"tracking-service/internal/domain"
	"tracking-service/internal/service"
)

// TrackingHandler menangani semua HTTP request untuk tracking
type TrackingHandler struct {
	service *service.TrackingService
}

// NewTrackingHandler membuat instance handler baru
func NewTrackingHandler(svc *service.TrackingService) *TrackingHandler {
	return &TrackingHandler{service: svc}
}

// writeJSON adalah helper untuk menulis response JSON
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError adalah helper untuk response error
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// =========================================================
// POST /api/v1/tracking/events
// Mencatat event tracking baru (dipanggil dari internal service atau API)
// Body: { "awb": "...", "status": "...", "hub_id": "...", "location": "...", "timestamp": "..." }
// =========================================================

func (h *TrackingHandler) RecordEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req domain.AddTrackingEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	event, err := h.service.RecordEvent(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, event)
}

// =========================================================
// GET /api/v1/tracking/{awb}/history
// Mengambil riwayat lengkap perjalanan paket (dari MongoDB)
// Response: { "awb": "...", "events": [...], "total": N }
// =========================================================

func (h *TrackingHandler) GetTrackingHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract AWB dari URL path: /api/v1/tracking/{awb}/history
	awb := extractAWBFromPath(r.URL.Path, "/history")
	if awb == "" {
		writeError(w, http.StatusBadRequest, "AWB tidak ditemukan di URL")
		return
	}

	history, err := h.service.GetTrackingHistory(r.Context(), awb)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, history)
}

// =========================================================
// GET /api/v1/tracking/{awb}/status
// Mengambil status TERAKHIR paket (dari Redis jika ada, fallback ke MongoDB)
// Response: { "awb": "...", "current_status": "...", "last_location": "...", "last_updated": "..." }
// =========================================================

func (h *TrackingHandler) GetCurrentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	awb := extractAWBFromPath(r.URL.Path, "/status")
	if awb == "" {
		writeError(w, http.StatusBadRequest, "AWB tidak ditemukan di URL")
		return
	}

	status, err := h.service.GetCurrentStatus(r.Context(), awb)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, status)
}

// =========================================================
// GET /health — Liveness probe untuk Kubernetes
// GET /ready  — Readiness probe untuk Kubernetes
// =========================================================

func (h *TrackingHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "tracking-service",
	})
}

// =========================================================
// HELPER
// =========================================================

// extractAWBFromPath mengekstrak AWB dari URL path.
// Contoh: "/api/v1/tracking/JKT-2024-001/history" → "JKT-2024-001"
func extractAWBFromPath(path string, suffix string) string {
	// Hapus suffix dulu
	path = strings.TrimSuffix(path, suffix)
	// Ambil segment terakhir
	parts := strings.Split(strings.TrimRight(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}