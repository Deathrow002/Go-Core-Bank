package main

import (
	"account-service/internal/account/controllers"
	"account-service/internal/account/repository"
	"account-service/internal/account/service"
	"account-service/internal/config"
	"account-service/internal/database"
	"account-service/internal/external"
	"account-service/internal/router" // Import router package
	"log"
	"os"
)

// @title Core Banking Account Service API
// @version 1.0
// @description A microservice for managing bank accounts
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run database migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize dependencies
	db := database.GetDB()
	accountRepo := repository.NewAccountRepository(db)

	customerServiceURL := os.Getenv("CUSTOMER_SERVICE_URL")
	if customerServiceURL == "" {
		customerServiceURL = "http://localhost:8080"
	}
	customerClient := external.NewCustomerClient(customerServiceURL)

	accountService := service.NewAccountService(accountRepo, customerClient)
	accountController := controllers.NewAccountController(accountService)

	// 🎯 **USE SEPARATE ROUTER**
	r := router.SetupRouter(cfg, accountController)

	// Start server
	log.Printf("Starting Account Service on %s", cfg.GetServerAddress())
	if err := r.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}