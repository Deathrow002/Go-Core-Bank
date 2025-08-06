package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "authentication-service",
			"version": "1.0.0",
		})
	})

	// Add HEAD /health route for health checks
	router.HEAD("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ready",
		})
	})
}