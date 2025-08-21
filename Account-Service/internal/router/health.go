package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// setupHealthRoutes configures health check routes
// Custom Prometheus metrics
var (
    accountCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "account_created_total",
        Help: "Total number of accounts created.",
    })
    accountUpdateTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "account_update_total",
        Help: "Total number of account updates.",
    })
    accountPasswordResetTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "account_password_reset_total",
        Help: "Total number of password resets.",
    })
    accountErrorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
        Name: "account_error_total",
        Help: "Total number of account errors by code.",
    }, []string{"code"})
)

func init() {
    prometheus.MustRegister(accountCreatedTotal)
    prometheus.MustRegister(accountUpdateTotal)
    prometheus.MustRegister(accountPasswordResetTotal)
    prometheus.MustRegister(accountErrorTotal)
}

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

    // Demo: increment metrics on /health for visibility
    router.GET("/service-metrics", func(c *gin.Context) {
        accountCreatedTotal.Inc()
        accountUpdateTotal.Inc()
        accountPasswordResetTotal.Inc()
        accountErrorTotal.WithLabelValues("400").Inc()
        accountErrorTotal.WithLabelValues("500").Inc()
        c.JSON(http.StatusOK, gin.H{"message": "Demo metrics incremented."})
    })

    // Prometheus metrics endpoint
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}