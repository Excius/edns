package middleware

import (
	"net/http"

	"github.com/excius/edns/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {

		logger.FromContext(c.Request.Context()).Error(
			"panic recovered",
			zap.Any("panic", recovered),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal server error",
			},
		)
	})
}
