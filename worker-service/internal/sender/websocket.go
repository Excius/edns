package sender

import (
	"context"

	"github.com/excius/edns/internal/logger"
	"go.uber.org/zap"
)

type WebSocketSender struct{}

func NewWebSocketSender() *WebSocketSender {
	return &WebSocketSender{}
}

func (w *WebSocketSender) Send(ctx context.Context, notificationID string) error {
	logger.Log.Info(
		"WebSocket message sent",
		zap.String("notification_id", notificationID),
	)

	return nil
}
