package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Custom Prometheus metrics
var (
	customerProfileQueryTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "customer_profile_query_total",
		Help: "Total number of customer profile queries.",
	})
	customerProfileUpdateTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "customer_profile_update_total",
		Help: "Total number of customer profile updates.",
	})
	customerKYCUpdateTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "customer_kyc_update_total",
		Help: "Total number of KYC updates.",
	})
	customerLatencyBucket = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "customer_latency_bucket",
		Help: "Latency distribution for customer operations.",
		Buckets: prometheus.DefBuckets,
	})
	customerErrorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "customer_error_total",
		Help: "Total number of customer errors by endpoint.",
	}, []string{"endpoint"})
	customerActivityTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "customer_activity_total",
		Help: "Total customer activity by customer_id.",
	}, []string{"customer_id"})
	customerApiCallsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "customer_api_calls_total",
		Help: "Total number of API calls to customer service.",
	})
)

func init() {
	prometheus.MustRegister(customerProfileQueryTotal)
	prometheus.MustRegister(customerProfileUpdateTotal)
	prometheus.MustRegister(customerKYCUpdateTotal)
	prometheus.MustRegister(customerLatencyBucket)
	prometheus.MustRegister(customerErrorTotal)
	prometheus.MustRegister(customerActivityTotal)
	prometheus.MustRegister(customerApiCallsTotal)
}

func setupHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "customer-service",
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
		customerProfileQueryTotal.Inc()
		customerProfileUpdateTotal.Inc()
		customerKYCUpdateTotal.Inc()
		customerLatencyBucket.Observe(0.3)
		customerErrorTotal.WithLabelValues("/profile").Inc()
		customerErrorTotal.WithLabelValues("/kyc").Inc()
		customerActivityTotal.WithLabelValues("123").Inc()
		customerApiCallsTotal.Inc()
		c.JSON(http.StatusOK, gin.H{"message": "Demo metrics incremented."})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}