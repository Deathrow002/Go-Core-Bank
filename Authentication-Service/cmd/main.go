package main

import (
	"authentication-service/internal/authentication/controller"
	"authentication-service/internal/authentication/repository"
	"authentication-service/internal/authentication/service"
	"authentication-service/internal/config"
	"authentication-service/internal/database"
	"authentication-service/internal/router"
	"log"
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

	// Initialize dependencies
	db := database.GetDB()
	authenticationRepo := repository.NewAuthenticationRepository(db)
	authenticationService := service.NewAuthenticationService(authenticationRepo)
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