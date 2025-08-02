package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yourname/employee-api/config"
	"github.com/yourname/employee-api/middleware"
	"github.com/yourname/employee-api/utils"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Message  string `json:"message"`
}

// HealthCheckHandler handles health check requests
func HealthCheckHandler(c *gin.Context) {
	logger := utils.GetLogger()
	requestID := middleware.GetRequestID(c)

	// Log health check request
	logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"operation":  "health_check",
	}).Info("Processing health check request")

	response := HealthResponse{
		Status:  "ok",
		Message: "GormTest application is running",
	}

	// Check database connectivity
	if err := config.HealthCheck(); err != nil {
		// Log database health check failure
		utils.LogDBError(c, "health_check", err)

		response.Status = "degraded"
		response.Database = "unhealthy"
		response.Message = "Application is running but database is unavailable"

		// Log degraded health status
		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"operation":  "health_check",
			"status":     "degraded",
		}).Warn("Health check returned degraded status")

		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	response.Database = "healthy"

	// Log successful health check
	logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"operation":  "health_check",
		"status":     "healthy",
	}).Info("Health check completed successfully")

	c.JSON(http.StatusOK, response)
}
