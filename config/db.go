package config

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB initializes the database connection
func InitDB() error {
	// Get connection string from environment variable or use default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://admin:1234@localhost:5439/appdb"
	}

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitDB() first.")
	}
	return db
}

// CloseDB closes the database connection gracefully
func CloseDB() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	log.Println("Closing database connection...")
	return sqlDB.Close()
}

// GracefulShutdown handles graceful shutdown of database connections
func GracefulShutdown(ctx context.Context) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Create a channel to signal completion
	done := make(chan error, 1)

	go func() {
		log.Println("Gracefully shutting down database connection...")
		done <- sqlDB.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("Error during database shutdown: %v", err)
			return err
		}
		log.Println("Database connection closed successfully")
		return nil
	case <-ctx.Done():
		log.Println("Database shutdown timed out")
		return ctx.Err()
	}
}

// HealthCheck performs a simple health check on the database
func HealthCheck() error {
	if db == nil {
		return gorm.ErrInvalidDB
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}
