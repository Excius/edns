package httpx

import (
	"errors"
	"net/http"

	"github.com/excius/edns/internal/apperrors"
	"github.com/excius/edns/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})

	case errors.Is(err, apperrors.ErrNotificationNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "notification not found",
		})

	case errors.Is(err, apperrors.ErrDeliveryNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": "delivery not found",
		})

	default:
		logger.Log.Error(
			"request failed",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}
