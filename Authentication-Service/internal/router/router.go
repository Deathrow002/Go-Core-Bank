package router

import (
	"authentication-service/internal/authentication/controller"
	"authentication-service/internal/config"
	"authentication-service/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the main router
func SetupRouter(cfg *config.Config, authenticationController *controller.AuthenticationController) *gin.Engine {
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
	setupAuthenticationRoutes(v1, authenticationController)

	return router
}