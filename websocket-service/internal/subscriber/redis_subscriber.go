package subscriber

import (
	"context"
	"fmt"

	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/websocket-service/internal/metrics"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisSubscriber struct {
	client  *redis.Client
	channel string
	metrics *metrics.Metrics
}

func NewRedisSubscriber(client *redis.Client, channel string, metrics *metrics.Metrics) *RedisSubscriber {
	return &RedisSubscriber{
		client:  client,
		channel: channel,
		metrics: metrics,
	}
}

func (s *RedisSubscriber) Start(ctx context.Context, handler EventHandler) error {

	pubsub := s.client.Subscribe(ctx, s.channel)
	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("pubsub receive failed: %w", err)
	}

	logger.Log.Info(
		"Subscribed to Redis channel",
		zap.String("channel", s.channel),
	)

	messages := pubsub.Channel()

	for {
		select {

		case <-ctx.Done():
			return nil

		case msg, ok := <-messages:
			if !ok {
				return fmt.Errorf("publish channel closed")
			}

			logger.Log.Info(
				"Message received from redis",
				zap.String("payload", msg.Payload),
			)

			s.metrics.Subscriber.MessagesReceived.Inc()

			if err := handler.Handle(ctx, []byte(msg.Payload)); err != nil {
				s.metrics.Subscriber.MessageHandlingErrors.Inc()
				logger.Log.Error(
					"message handler failed",
					zap.Error(err),
				)
			}
		}
	}
}
