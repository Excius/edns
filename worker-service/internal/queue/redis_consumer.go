package queue

import (
	"context"
	"strings"

	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/internal/stream"
	"github.com/excius/edns/worker-service/internal/metrics"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisConsumer struct {
	clinet       *redis.Client
	stream       string
	group        string
	consumerName string
	metrics      *metrics.ConsumerMetrics
}

func NewRedisConsumer(client *redis.Client, stream string, group string, consumerName string, metrics *metrics.ConsumerMetrics) *RedisConsumer {
	return &RedisConsumer{
		clinet:       client,
		stream:       stream,
		group:        group,
		consumerName: consumerName,
		metrics:      metrics,
	}
}

func (r *RedisConsumer) EnsureGroup(ctx context.Context) error {
	err := r.clinet.XGroupCreateMkStream(ctx, r.stream, r.group, "$").Err()
	if err != nil {
		if strings.Contains(err.Error(), "BUSYGROUP") {
			return nil
		}
		return err
	}
	return nil
}

func (r *RedisConsumer) Start(ctx context.Context, processor stream.Processor) error {

	for {
		streams, err := r.clinet.XReadGroup(ctx, &redis.XReadGroupArgs{
			Streams:  []string{r.stream, ">"},
			Group:    r.group,
			Consumer: r.consumerName,
			Block:    0,
			Count:    1,
		}).Result()

		if err != nil {
			return err
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {

				r.metrics.MessagesReceived.Inc()

				// Before processing each message
				select {
				case <-ctx.Done():
					return nil
				default:
				}

				result, err := processor.Process(ctx, msg.Values)
				if err != nil {
					r.metrics.MessageProcessingErrors.Inc()
					logger.Log.Error(
						"failed processing message",
						zap.String("message_id", msg.ID),
						zap.Error(err),
					)

					// Leave message in PEL
					continue
				}

				if !result.Ack {
					logger.Log.Debug(
						"notification still in progress",
						zap.String("message_id", msg.ID),
					)

					// Leave message in PEL
					continue
				}

				acked, err := r.clinet.XAck(ctx, r.stream, r.group, msg.ID).Result()
				if err != nil {
					logger.Log.Error(
						"failed to acknowledge message",
						zap.String("message_id", msg.ID),
						zap.Error(err),
					)

					continue
				}

				r.metrics.MessagesAcknowledged.Inc()

				logger.Log.Debug(
					"message acknowledged",
					zap.String("message_id", msg.ID),
					zap.Int64("acked", acked),
				)
			}
		}
	}
}
