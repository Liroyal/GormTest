package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yourname/employee-api/config"
	"github.com/yourname/employee-api/handlers"
	"github.com/yourname/employee-api/middleware"
	"github.com/yourname/employee-api/utils"
)

func main() {
	// Initialize centralized logger
	logger := utils.InitLogger()
	logger.Info("GormTest application starting...")

	// Load configuration
	cfg := config.Load()
	logger.WithFields(logrus.Fields{
		"server_port": cfg.Server.Port,
		"server_host": cfg.Server.Host,
	}).Info("Configuration loaded")

	// Initialize database connection
	if err := config.InitDB(); err != nil {
		logger.WithError(err).Fatal("Failed to initialize database connection")
	}
	logger.Info("Database connection initialized successfully")

	// Create Gin router with centralized middleware
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.ErrorHandler(logger))

	// Register routes
	registerRoutes(router)
	logger.Info("Routes registered successfully")

	// Configure server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.WithField("address", serverAddr).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	// Close database connection
	if err := config.CloseDB(); err != nil {
		logger.WithError(err).Error("Error closing database connection")
	}

	logger.Info("Server exited gracefully")
}

// registerRoutes registers all application routes
func registerRoutes(router *gin.Engine) {
	// Health check route
	router.GET("/health", handlers.HealthCheckHandler)

	// Employee routes
	router.POST("/employees", handlers.CreateEmployeeHandler)
	router.GET("/employees/:id", handlers.GetEmployeeHandler)
	router.PUT("/employees/:id", handlers.UpdateEmployeeHandler)
}
