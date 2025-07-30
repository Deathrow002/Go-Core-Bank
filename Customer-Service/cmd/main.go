package main

import (
	"customer-service/internal/config"
	"customer-service/internal/customer/controllers"
	"customer-service/internal/customer/repository"
	"customer-service/internal/customer/service"
	"customer-service/internal/database"
	"customer-service/internal/router" // Import the router package
	"log"
)

// @title Core Banking Customer Service API
// @version 1.0
// @description A microservice for managing bank customers
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
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
	customerRepo := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepo)
	customerController := controllers.NewCustomerController(customerService)

	// 🎯 **Setup router using separate router package**
	r := router.SetupRouter(cfg, customerController)

	// Debug: Print all registered routes
	routes := r.Routes()
	log.Println("📋 Registered routes:")
	for _, route := range routes {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	// Start server
	log.Printf("Starting Customer Service on %s", cfg.GetServerAddress())
	if err := r.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
