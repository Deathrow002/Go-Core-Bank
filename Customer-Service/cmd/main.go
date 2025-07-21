package main

import (
	"customer-service/internal/config"
	"customer-service/internal/customer/controllers"
	"customer-service/internal/customer/repository"
	"customer-service/internal/customer/service"
	"customer-service/internal/database"
	"customer-service/pkg/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// Setup router
	router := setupRouter(cfg, customerController)

	// Start server
	log.Printf("Starting server on %s", cfg.GetServerAddress())
	if err := router.Run(cfg.GetServerAddress()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(cfg *config.Config, customerController *controllers.CustomerController) *gin.Engine {
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
			"status": "healthy",
			"service": "customer-service",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		customers := v1.Group("/customers")
		{
			customers.POST("", customerController.CreateCustomer)
			customers.GET("/:id", customerController.GetCustomer)
			customers.PUT("/:id", customerController.UpdateCustomer)
			customers.DELETE("/:id", customerController.DeleteCustomer)
			customers.GET("", customerController.ListCustomers)
			customers.GET("/search", customerController.SearchCustomers)
		}
	}

	return router
}
