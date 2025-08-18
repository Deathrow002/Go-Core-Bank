package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

    // Add HEAD /health route for health checks
	router.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

    router.GET("/ready", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ready",
        })
    })

    // Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}