// Package health provides HTTP handlers for health check endpoints.
package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// Status represents the health status of a component.
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

// ComponentHealth represents the health of a single component.
type ComponentHealth struct {
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
}

// HealthResponse represents the overall health status response.
type HealthResponse struct {
	Status     Status                     `json:"status"`
	Components map[string]ComponentHealth `json:"components"`
	Timestamp  string                     `json:"timestamp"`
}

// Handler handles health check requests.
type Handler struct {
	db *sql.DB
}

// NewHandler creates a new health check handler.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// Handle processes health check requests.
// @Summary Health check
// @Description Check the health status of the API and its dependencies
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:     StatusHealthy,
		Components: make(map[string]ComponentHealth),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	// Check database connectivity
	dbHealth := h.checkDatabase()
	response.Components["database"] = dbHealth
	if dbHealth.Status == StatusUnhealthy {
		response.Status = StatusUnhealthy
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Status == StatusUnhealthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}

// checkDatabase verifies database connectivity.
func (h *Handler) checkDatabase() ComponentHealth {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		return ComponentHealth{
			Status:  StatusUnhealthy,
			Message: "database connection failed",
		}
	}

	return ComponentHealth{
		Status:  StatusHealthy,
		Message: "connected",
	}
}
