package router

import (
	"transaction-service/internal/config"
	"transaction-service/internal/transaction/controller"
	"transaction-service/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the main router
func SetupRouter(cfg *config.Config, transactionController *controller.TransactionController) *gin.Engine {
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
	setupTrasactionRoutes(v1, transactionController)

	return router
}