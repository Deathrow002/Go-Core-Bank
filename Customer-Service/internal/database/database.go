package database

import (
	"customer-service/internal/config"
	"customer-service/internal/customer/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase(cfg *config.Config) error {
	return initDatabaseWithRetry(cfg, 10, 5*time.Second)
}

// initDatabaseWithRetry initializes the database connection with retry logic
func initDatabaseWithRetry(cfg *config.Config, maxRetries int, retryDelay time.Duration) error {
	var err error

	// Configure GORM logger
	var gormLogger logger.Interface
	if cfg.IsDevelopment() {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Try to connect with retries
	for i := 0; i < maxRetries; i++ {
		// Connect to database
		DB, err = gorm.Open(postgres.Open(cfg.GetDatabaseDSN()), &gorm.Config{
			Logger: gormLogger,
		})
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to connect to database: %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
		}

		// Test connection
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to get database instance: %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf("failed to get database instance after %d attempts: %w", maxRetries, err)
		}

		if err := sqlDB.Ping(); err != nil {
			log.Printf("Attempt %d/%d: Failed to ping database: %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
		}

		log.Println("Successfully connected to database")
		return nil
	}

	return fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Run auto-migration for all models
	err := DB.AutoMigrate(
		&models.Customer{},
	)
	if err != nil {
		return fmt.Errorf("failed to run auto-migration: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	return sqlDB.Close()
}
