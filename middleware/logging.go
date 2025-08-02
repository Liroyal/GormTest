package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Logger returns a gin.HandlerFunc (middleware) that logs requests using logrus.
func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Generate request ID if not present
		requestID := param.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Log with structured fields
		logger.WithFields(logrus.Fields{
			"request_id":    requestID,
			"timestamp":     param.TimeStamp.Format(time.RFC3339),
			"status":        param.StatusCode,
			"latency":       param.Latency,
			"client_ip":     param.ClientIP,
			"method":        param.Method,
			"path":          param.Path,
			"user_agent":    param.Request.UserAgent(),
			"error_message": param.ErrorMessage,
		}).Info("HTTP Request")

		return ""
	})
}

// RequestID is a middleware that generates or extracts request IDs
func RequestID() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set the request ID in both the response header and context
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	})
}

// GetRequestID extracts the request ID from the context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}
