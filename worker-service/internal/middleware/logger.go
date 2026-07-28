package middleware

import (
	"time"

	"github.com/excius/edns/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		latency := time.Since(start)

		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}

		logger.FromContext(c.Request.Context()).Info(
			"http request",
			zap.String("method", c.Request.Method),
			zap.String("route", route),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}
