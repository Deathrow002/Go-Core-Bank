package main

import (
	"account-service/internal/account/controllers"
	"account-service/internal/account/repository"
	"account-service/internal/account/service"
	"account-service/internal/config"
	"account-service/internal/database"
	"account-service/internal/external"
	"account-service/internal/router" // Import router package
	"account-service/pkg/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

func setupRouter(cfg *config.Config, accountController *controllers.AccountController) *gin.Engine {
	// Set gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "account-service",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		accounts := v1.Group("/accounts")
		{
			accounts.POST("", accountController.CreateAccount)
			accounts.GET("/:id", accountController.GetAccount)
			accounts.PUT("/:id", accountController.UpdateAccount)
			accounts.DELETE("/:id", accountController.DeleteAccount)
			accounts.GET("", accountController.ListAccounts)
			accounts.GET("/search", accountController.SearchAccounts)
			accounts.GET("/number/:account_number", accountController.GetAccountByNumber)
			accounts.GET("/customer/:customer_id", accountController.GetAccountsByCustomer)
		}
	}

	return router
}