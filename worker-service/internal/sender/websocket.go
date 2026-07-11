package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/excius/edns/internal/events"
	"github.com/excius/edns/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type WebSocketSender struct {
	client  *redis.Client
	channel string
}

func NewWebSocketSender(client *redis.Client, channel string) *WebSocketSender {
	return &WebSocketSender{
		client:  client,
		channel: channel,
	}
}

func (w *WebSocketSender) Send(ctx context.Context, event events.NotificationEvent) error {

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("event marshel failed: %w", err)
	}

	if err := w.client.Publish(ctx, w.channel, payload).Err(); err != nil {
		return fmt.Errorf("publish payload failed: %w", err)
	}

	logger.Log.Info(
		"WebSocket message sent",
		zap.String("notification_id", event.NotificationID),
	)

	return nil
}
