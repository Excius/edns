package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// FromContext extracts a zap.Logger from the provided context.
// If a request ID is found in the context under RequestIDKey, it returns a logger with the "request_id" field attached.
// Otherwise, it returns the global Log instance.
func FromContext(ctx context.Context) *zap.Logger {
	if ctx != nil {
		if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
			return Log.With(zap.String("request_id", reqID))
		}
	}
	return Log
}
