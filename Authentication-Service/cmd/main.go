package main

import (
	"authentication-service/internal/authentication/controller"
	"authentication-service/internal/authentication/models"
	"authentication-service/internal/authentication/repository"
	"authentication-service/internal/authentication/service"
	"authentication-service/internal/config"
	"authentication-service/internal/database"
	"authentication-service/internal/external"
	"authentication-service/internal/router"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
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
	if err := initAdminAccount(); err != nil {
		log.Fatalf("Failed to initialize admin account: %v", err)
	}

	// Get the customer service URL from environment variable or default
	customerServiceURL := os.Getenv("CUSTOMER_SERVICE_URL")
	if customerServiceURL == "" {
		customerServiceURL = "http://localhost:8080"
	}

	// Initialize dependencies
	customerClient := external.NewCustomerClient(customerServiceURL)

	// Get the database instance
	db := database.GetDB()

	// Initialize authentication service with customer client
	authenticationRepo := repository.NewAuthenticationRepository(db)
	authenticationService := service.NewAuthenticationService(authenticationRepo, customerClient)
	authenticationController := controller.NewAuthenticationController(authenticationService)

	// Setup router
	r := router.SetupRouter(cfg, authenticationController)

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

// Add this function to the same file or import it if defined elsewhere
func initAdminAccount() error {
	db := database.GetDB()
	var count int64
	if err := db.Model(&models.Authentication{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("Admin account already exists, skipping admin initialization.")
		return nil
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminEmail := os.Getenv("ADMIN_EMAIL")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.Authentication{
		Username:     adminUsername,
		PasswordHash: string(hashedPassword),
		Email:        adminEmail,
		Role:         "admin",
		IsLocked:     false,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	log.Println("Default admin account created.")
	return nil
}

