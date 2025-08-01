package main

import (
	"authentication-service/internal/authentication/models"
	"authentication-service/internal/config"
	"authentication-service/internal/database"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
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

	// Now it's safe to initialize admin account
	if err := initAdminAccount(); err != nil {
		log.Fatalf("Failed to initialize admin account: %v", err)
	}

	log.Println("Migrations and admin account initialization completed successfully")
}

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

	// Use env vars or defaults
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
		Role:         string(models.RoleTypeAdmin),
		IsLocked:     false,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	log.Println("Default admin account created.")
	return nil
}