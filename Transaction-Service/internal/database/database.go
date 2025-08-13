package database

import (
	"fmt"
	"log"
	"time"
	"transaction-service/internal/config"
	"transaction-service/internal/transaction/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) error {
	return initDatabaseWithRetry(cfg, 10, 5*time.Second)
}

// initDatabaseWithRetry tries to initialize the database with retries.
func initDatabaseWithRetry(cfg *config.Config, maxRetries int, retryInterval time.Duration) error {
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
				time.Sleep(retryInterval)
				continue
			}
			return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
		}

		// Test connection
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to get database instance: %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(retryInterval)
				continue
			}
			return fmt.Errorf("failed to get database instance after %d attempts: %w", maxRetries, err)
		}

		if err := sqlDB.Ping(); err != nil {
			log.Printf("Attempt %d/%d: Failed to ping database: %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(retryInterval)
				continue
			}
			return fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
		}

		log.Println("Successfully connected to database")
		return nil
	}

	return fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Run migrations
	if err := DB.AutoMigrate(
		&models.Transaction{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Database migrations completed successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

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