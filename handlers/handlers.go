package handlers

import (
	"encoding/json"
	"net/http"
)

// Handler represents the main handler struct
type Handler struct {
	// TODO: Add dependencies like database, logger, etc.
}

// NewHandler creates a new handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
		"message": "GormTest application is running",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SetupRoutes sets up HTTP routes
func (h *Handler) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	// TODO: Add more routes
	return mux
}
