package middleware

import (
	"context"

	"github.com/excius/edns/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderXRequestID = "X-Request-ID"

// RequestID middleware ensures every request carries a unique request ID in context and response headers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Header(HeaderXRequestID, reqID)

		ctx := context.WithValue(c.Request.Context(), logger.RequestIDKey, reqID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
