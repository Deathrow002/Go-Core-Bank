package router

import (
	"account-service/internal/account/controllers"
	"account-service/internal/config"
	"account-service/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// setupAPIRoutes configures API routes
func setupAPIRoutes(router *gin.Engine, accountController *controllers.AccountController) {
    v1 := router.Group("/api/v1")
    {
        setupAccountRoutes(v1, accountController)
    }
}

// SetupRouter configures and returns the main router
func SetupRouter(cfg *config.Config, accountController *controllers.AccountController) *gin.Engine {
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

    // Setup routes
    setupHealthRoutes(router)
    setupAPIRoutes(router, accountController)

    return router
}