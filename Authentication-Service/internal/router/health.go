package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Custom Prometheus metrics
var (
	authLoginSuccessTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_login_success_total",
		Help: "Total number of successful logins.",
	})
	authLoginFailureTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_login_failure_total",
		Help: "Total number of failed logins.",
	})
	authTokenIssuedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_token_issued_total",
		Help: "Total number of tokens issued.",
	})
	authLatencySeconds = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "auth_latency_seconds",
		Help: "Authentication latency in seconds.",
		Buckets: prometheus.DefBuckets,
	})
	authErrorCodeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_error_code_total",
		Help: "Total number of authentication errors by code.",
	}, []string{"code"})
)

func init() {
	prometheus.MustRegister(authLoginSuccessTotal)
	prometheus.MustRegister(authLoginFailureTotal)
	prometheus.MustRegister(authTokenIssuedTotal)
	prometheus.MustRegister(authLatencySeconds)
	prometheus.MustRegister(authErrorCodeTotal)
}

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

	// Demo: increment metrics for visibility
	router.GET("/demo-metrics", func(c *gin.Context) {
		authLoginSuccessTotal.Inc()
		authLoginFailureTotal.Inc()
		authTokenIssuedTotal.Inc()
		authLatencySeconds.Observe(0.2)
		authErrorCodeTotal.WithLabelValues("401").Inc()
		authErrorCodeTotal.WithLabelValues("403").Inc()
		authErrorCodeTotal.WithLabelValues("500").Inc()
		c.JSON(http.StatusOK, gin.H{"message": "Demo metrics incremented."})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}