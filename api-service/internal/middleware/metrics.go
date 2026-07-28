package middleware

import (
	"strconv"
	"time"

	"github.com/excius/edns/api-service/internal/metrics"
	"github.com/gin-gonic/gin"
)

func Metrics(metrics *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {

		metrics.HTTP.RequestsInProgress.Inc()
		defer metrics.HTTP.RequestsInProgress.Dec()

		start := time.Now()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}

		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		metrics.HTTP.RequestDuration.
			WithLabelValues(method, route, status).
			Observe(time.Since(start).Seconds())

		metrics.HTTP.RequestsTotal.
			WithLabelValues(method, route, status).
			Inc()
	}
}
