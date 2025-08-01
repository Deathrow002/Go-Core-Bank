package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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