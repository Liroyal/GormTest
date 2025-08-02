package utils

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var AppLogger *logrus.Logger

// InitLogger initializes the application logger
func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logger.SetLevel(logrus.InfoLevel)

	AppLogger = logger
	return logger
}

// GetLogger returns the application logger instance
func GetLogger() *logrus.Logger {
	if AppLogger == nil {
		return InitLogger()
	}
	return AppLogger
}

// LogDBError logs database-related errors with context
func LogDBError(ctx context.Context, operation string, err error, additionalFields ...logrus.Fields) {
	logger := GetLogger()
	
	fields := logrus.Fields{
		"operation": operation,
		"error":     err.Error(),
		"type":      "database_error",
	}

	// Add request ID if available from Gin context
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if requestID, exists := ginCtx.Get("request_id"); exists {
			fields["request_id"] = requestID
		}
	}

	// Merge additional fields
	for _, additionalField := range additionalFields {
		for k, v := range additionalField {
			fields[k] = v
		}
	}

	// Check if it's a specific GORM error
	if err == gorm.ErrRecordNotFound {
		logger.WithFields(fields).Warn("Database record not found")
	} else {
		logger.WithFields(fields).Error("Database operation failed")
	}
}

// LogValidationError logs validation errors with context
func LogValidationError(ctx context.Context, field string, value interface{}, err error, additionalFields ...logrus.Fields) {
	logger := GetLogger()
	
	fields := logrus.Fields{
		"field":         field,
		"value":         value,
		"error":         err.Error(),
		"type":          "validation_error",
	}

	// Add request ID if available from Gin context
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if requestID, exists := ginCtx.Get("request_id"); exists {
			fields["request_id"] = requestID
		}
	}

	// Merge additional fields
	for _, additionalField := range additionalFields {
		for k, v := range additionalField {
			fields[k] = v
		}
	}

	logger.WithFields(fields).Warn("Validation error occurred")
}

// LogBusinessError logs business logic errors with context
func LogBusinessError(ctx context.Context, operation string, err error, additionalFields ...logrus.Fields) {
	logger := GetLogger()
	
	fields := logrus.Fields{
		"operation": operation,
		"error":     err.Error(),
		"type":      "business_error",
	}

	// Add request ID if available from Gin context
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if requestID, exists := ginCtx.Get("request_id"); exists {
			fields["request_id"] = requestID
		}
	}

	// Merge additional fields
	for _, additionalField := range additionalFields {
		for k, v := range additionalField {
			fields[k] = v
		}
	}

	logger.WithFields(fields).Error("Business logic error occurred")
}

// LogInfo logs informational messages with context
func LogInfo(ctx context.Context, message string, additionalFields ...logrus.Fields) {
	logger := GetLogger()
	
	fields := logrus.Fields{
		"message": message,
	}

	// Add request ID if available from Gin context
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if requestID, exists := ginCtx.Get("request_id"); exists {
			fields["request_id"] = requestID
		}
	}

	// Merge additional fields
	for _, additionalField := range additionalFields {
		for k, v := range additionalField {
			fields[k] = v
		}
	}

	logger.WithFields(fields).Info(message)
}

// WithRequestID is a helper to create a logger with request ID context
func WithRequestID(c *gin.Context) *logrus.Entry {
	logger := GetLogger()
	requestID := ""
	
	if c != nil {
		if id, exists := c.Get("request_id"); exists {
			requestID = id.(string)
		}
	}

	return logger.WithField("request_id", requestID)
}
