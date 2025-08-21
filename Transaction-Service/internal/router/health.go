package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Custom Prometheus metrics
var (
	transactionTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "transaction_total",
		Help: "Total number of transactions.",
	})
	transactionValueTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "transaction_value_total",
		Help: "Total value of transactions.",
	})
	transactionSuccessTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "transaction_success_total",
		Help: "Total number of successful transactions.",
	})
	transactionFailureTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "transaction_failure_total",
		Help: "Total number of failed transactions by reason.",
	}, []string{"reason"})
	transactionProcessingTimeSeconds = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "transaction_processing_time_seconds",
		Help: "Transaction processing time in seconds.",
		Buckets: prometheus.DefBuckets,
	})
	transactionTopVolume = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "transaction_top_volume",
		Help: "Top customers/regions by transaction volume.",
	}, []string{"customer_id", "region"})
)

func init() {
	prometheus.MustRegister(transactionTotal)
	prometheus.MustRegister(transactionValueTotal)
	prometheus.MustRegister(transactionSuccessTotal)
	prometheus.MustRegister(transactionFailureTotal)
	prometheus.MustRegister(transactionProcessingTimeSeconds)
	prometheus.MustRegister(transactionTopVolume)
}

func setupHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "transaction-service",
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
	router.GET("/service-metrics", func(c *gin.Context) {
		transactionTotal.Inc()
		transactionValueTotal.Add(100)
		transactionSuccessTotal.Inc()
		transactionFailureTotal.WithLabelValues("insufficient_funds").Inc()
		transactionFailureTotal.WithLabelValues("invalid_account").Inc()
		transactionFailureTotal.WithLabelValues("system_error").Inc()
		transactionProcessingTimeSeconds.Observe(0.5)
		transactionTopVolume.WithLabelValues("123", "Bangkok").Inc()
		c.JSON(http.StatusOK, gin.H{"message": "Demo metrics incremented."})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}