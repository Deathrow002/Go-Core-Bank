package router

import (
	"account-service/internal/account/controllers"
	"account-service/internal/config"
	"account-service/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

// setupHealthRoutes configures health check routes
func setupHealthRoutes(router *gin.Engine) {
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":    "healthy",
            "service":   "account-service",
            "version":   "1.0.0",
            "timestamp": gin.H{},
        })
    })

    router.GET("/ready", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ready",
        })
    })
}

// setupAPIRoutes configures API routes
func setupAPIRoutes(router *gin.Engine, accountController *controllers.AccountController) {
    v1 := router.Group("/api/v1")
    {
        setupAccountRoutes(v1, accountController)
    }
}

// setupAccountRoutes configures account-specific routes
func setupAccountRoutes(rg *gin.RouterGroup, accountController *controllers.AccountController) {
    accounts := rg.Group("/accounts")
    {
        // Basic CRUD operations
        accounts.POST("", accountController.CreateAccount)
        accounts.GET("/:id", accountController.GetAccount)
        accounts.PUT("/:id", accountController.UpdateAccount)
        accounts.DELETE("/:id", accountController.DeleteAccount)

        // Listing and search
        accounts.GET("", accountController.ListAccounts)
        accounts.GET("/search", accountController.SearchAccounts)

        // Special queries
        accounts.GET("/number/:account_number", accountController.GetAccountByNumber)
        accounts.GET("/customer/:customer_id", accountController.GetAccountsByCustomer)
    }
}