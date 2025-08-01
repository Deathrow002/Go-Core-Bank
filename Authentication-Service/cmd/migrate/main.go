package main

import (
	"authentication-service/internal/config"
	"authentication-service/internal/database"
	"log"
)

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

	// Run migrations BEFORE initializing admin account
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations and admin account initialization completed successfully")
}