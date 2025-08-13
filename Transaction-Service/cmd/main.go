package main

import (
	"log"
	"os"
	"transaction-service/internal/config"
	"transaction-service/internal/database"
	external "transaction-service/internal/external"
	producer "transaction-service/internal/external/producer"
	"transaction-service/internal/router"
	"transaction-service/internal/transaction/controller"
	"transaction-service/internal/transaction/repository"
	"transaction-service/internal/transaction/service"

	"github.com/segmentio/kafka-go"
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

	accountServiceURL := os.Getenv("ACCOUNT_SERVICE_URL")
	if accountServiceURL == "" {
		accountServiceURL = "http://localhost:8082"
	}

	// Initialize dependencies
	accountClient := external.NewAccountClient(accountServiceURL)

	// Get the database instance
	db := database.GetDB()

	// Initialize account service with customer client
	transactionRepository := repository.NewTransactionRepository(db)

	// Kafka writer setup
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "account-balance-update",
		Balancer: &kafka.LeastBytes{},
	}

	transactionService := service.NewTransactionService(transactionRepository, accountClient, kafkaWriter)
	transactionController := controller.NewTransactionController(transactionService)

	// 🎯 **USE SEPARATE ROUTER**
	r := router.SetupRouter(cfg, transactionController)

	// Debug: Print all registered routes
	routes := r.Routes()
	log.Println("📋 Registered routes:")
	for _, route := range routes {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	r.POST("/update-balance", producer.SendBalanceUpdateToKafka(kafkaWriter))

	// Start server
	log.Printf("Starting Transaction Service on %s", cfg.GetServerAddress())
	if err := r.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}