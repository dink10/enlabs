package server

import (
	"net/http"
)

// HealthCheck provides health-check functionality.
// @Summary Health check
// @Description End-point providing health-check functionality
// @ID health-check
// @Tags Health Check
// @Produce json
// @Success 200 {object} server.Response "Status"
// @Router /health [get]
func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RenderResponse(w, r, NewResponse(http.StatusOK))
	}
}
