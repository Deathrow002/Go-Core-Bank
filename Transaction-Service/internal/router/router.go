package router

import (
	"transaction-service/internal/config"
	"transaction-service/internal/transaction/controller"
	"transaction-service/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// setupAPIRoutes configures API routes
func setupAPIRoutes(router *gin.Engine, accountController *controller.TransactionController) {
    v1 := router.Group("/api/v1")
    {
        setupTrasactionRoutes(v1, accountController)
    }
}

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

	// Setup customer routes
    setupAPIRoutes(router, transactionController)

	return router
}