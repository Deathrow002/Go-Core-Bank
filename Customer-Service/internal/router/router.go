package router

import (
	"customer-service/internal/config"
	"customer-service/internal/customer/controllers"
	"customer-service/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, customerController *controllers.CustomerController) *gin.Engine {
	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Setup health routes
	setupHealthRoutes(router)

	// Setup API routes
	api := router.Group("/api")
	v1 := api.Group("/v1")

	// Setup customer routes
	setupCustomerRoutes(v1, customerController)

	return router
}