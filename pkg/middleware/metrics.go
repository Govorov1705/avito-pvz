package middleware

import (
	"github.com/Govorov1705/avito-pvz/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics(c *gin.Context) {
	timer := prometheus.NewTimer(metrics.HttpRequestDuration)
	defer timer.ObserveDuration()

	metrics.HttpRequestsTotal.Inc()

	c.Next()
}
