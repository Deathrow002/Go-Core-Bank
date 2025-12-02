package main

import (
	"account-service/internal/account/controllers"
	"account-service/internal/account/repository"
	"account-service/internal/account/service"
	"account-service/internal/config"
	"account-service/internal/database"
	"account-service/internal/external"
	"account-service/internal/external/consumer"
	"account-service/internal/router"
	"context"
	"log"
	"os"
)

func main() {
	// This is the entry point for the Customer Service application.
	// The main function initializes the service, sets up the database,
	// and starts the HTTP server with the configured routes.

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := database.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// You should have this line to run migrations:
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize admin account after migrations
	customerServiceURL := os.Getenv("CUSTOMER_SERVICE_URL")
	if customerServiceURL == "" {
		customerServiceURL = "http://localhost:8080"
	}

	// Initialize dependencies
	customerClient := external.NewCustomerClient(customerServiceURL)

	// Get the database instance
	db := database.GetDB()

	// Initialize account service with customer client
	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo, customerClient)
	accountController := controllers.NewAccountController(accountService)

	// 🎯 **USE SEPARATE ROUTER**
	r := router.SetupRouter(cfg, accountController)

	// Debug: Print all registered routes
	routes := r.Routes()
	log.Println("📋 Registered routes:")
	for _, route := range routes {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	// Kafka consumer setup (use full cluster DNS as requested)
	brokers := []string{"kafka.core-bank.svc.cluster.local:9092"}
	topic := "account-balance-update"
	groupID := "account-balance-consumer-fixed" // Use fixed group ID

	consumer, err := consumer.NewBalanceUpdateConsumer(
		brokers,
		topic,
		groupID,
		accountService,
	)
	if err != nil {
		log.Fatalf("Failed to create balance update consumer: %v", err)
	}

	// Start consuming messages
	consumer.Start(context.Background())

	// Ensure the consumer stops gracefully
	defer consumer.Stop()

	// Start server
	log.Printf("Starting Account Service on %s", cfg.GetServerAddress())
	if err := r.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
