package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ErrorHandler is a middleware that handles errors and provides structured error responses
func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			requestID := GetRequestID(c)
			
			// Get the last error (most recent)
			err := c.Errors.Last()
			
			// Log the error with context
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Request.URL.Path,
				"method":     c.Request.Method,
				"error":      err.Error(),
				"error_type": err.Type,
			}).Error("Request error occurred")

			// If response hasn't been written yet, send error response
			if !c.Writer.Written() {
				var statusCode int
				var errorMessage string

				// Determine status code based on error type
				switch err.Type {
				case gin.ErrorTypeBind:
					statusCode = http.StatusBadRequest
					errorMessage = "Invalid request format"
				case gin.ErrorTypePublic:
					statusCode = http.StatusBadRequest
					errorMessage = err.Error()
				default:
					statusCode = http.StatusInternalServerError
					errorMessage = "Internal server error"
				}

				c.JSON(statusCode, ErrorResponse{
					Error: errorMessage,
				})
			}
		}
	})
}

// AbortWithError is a helper function to abort request with error and log it
func AbortWithError(c *gin.Context, statusCode int, err error, logger *logrus.Logger) {
	requestID := GetRequestID(c)
	
	// Log the error
	logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"path":       c.Request.URL.Path,
		"method":     c.Request.Method,
		"status":     statusCode,
		"error":      err.Error(),
	}).Error("Request aborted with error")

	// Return standardized error response
	c.JSON(statusCode, ErrorResponse{
		Error: err.Error(),
	})
	c.Abort()
}
