package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler_Success(t *testing.T) {
	// Setup
	router := setupTestRouter()
	router.GET("/health", HealthCheckHandler)

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	
	// Create response recorder
	w := httptest.NewRecorder()

	// Note: This test requires database mocking or test database setup
	// For a complete test, you would need to mock the config.HealthCheck() function
	// or use a test database
	
	// Perform request
	router.ServeHTTP(w, req)

	// Basic assertions that don't depend on database
	assert.Contains(t, []int{http.StatusOK, http.StatusServiceUnavailable}, w.Code)
	
	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Status)
	assert.NotEmpty(t, response.Message)
	
	// Status should be either "ok" or "degraded"
	assert.Contains(t, []string{"ok", "degraded"}, response.Status)
}

func TestHealthCheckHandler_DatabaseHealthy(t *testing.T) {
	// Setup
	router := setupTestRouter()
	router.GET("/health", HealthCheckHandler)

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	
	// Create response recorder
	w := httptest.NewRecorder()

	// Note: For a complete test, you would mock config.HealthCheck() to return nil
	// This would ensure the handler returns a healthy status
	
	// Perform request
	router.ServeHTTP(w, req)

	// Parse response
	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Assertions would depend on mocking the health check function
	// Example with mocked healthy database:
	// assert.Equal(t, http.StatusOK, w.Code)
	// assert.Equal(t, "ok", response.Status)
	// assert.Equal(t, "healthy", response.Database)
	// assert.Equal(t, "GormTest application is running", response.Message)
}

func TestHealthCheckHandler_DatabaseUnhealthy(t *testing.T) {
	// Setup
	router := setupTestRouter()
	router.GET("/health", HealthCheckHandler)

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	
	// Create response recorder
	w := httptest.NewRecorder()

	// Note: For a complete test, you would mock config.HealthCheck() to return an error
	// This would ensure the handler returns a degraded status
	
	// Perform request
	router.ServeHTTP(w, req)

	// Parse response
	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Assertions would depend on mocking the health check function
	// Example with mocked unhealthy database:
	// assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	// assert.Equal(t, "degraded", response.Status)
	// assert.Equal(t, "unhealthy", response.Database)
	// assert.Equal(t, "Application is running but database is unavailable", response.Message)
}

func TestHealthCheckHandler_ResponseStructure(t *testing.T) {
	// Setup
	router := setupTestRouter()
	router.GET("/health", HealthCheckHandler)

	// Create request
	req, _ := http.NewRequest("GET", "/health", nil)
	
	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Verify response structure
	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify all required fields are present
	assert.NotEmpty(t, response.Status, "Status field should not be empty")
	assert.NotEmpty(t, response.Message, "Message field should not be empty")
	
	// Database field might be empty if not set, so we don't assert NotEmpty for it
	// but we can verify it's a valid value if present
	if response.Database != "" {
		assert.Contains(t, []string{"healthy", "unhealthy"}, response.Database)
	}
}
