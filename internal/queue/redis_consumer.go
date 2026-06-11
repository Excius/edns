package queue

import (
	"context"
	"strings"

	"github.com/excius/edns/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisConsumer struct {
	clinet       *redis.Client
	stream       string
	group        string
	consumerName string
}

func NewRedisConsumer(client *redis.Client, stream string, group string, consumerName string) *RedisConsumer {
	return &RedisConsumer{
		clinet:       client,
		stream:       stream,
		group:        group,
		consumerName: consumerName,
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

func (r *RedisConsumer) Start(ctx context.Context, processor Processor) error {

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

				if err := processor.Process(ctx, msg.Values); err != nil {

					logger.Log.Error(
						"Failed processing message",
						zap.String("message_id", msg.ID),
						zap.Error(err),
					)
					continue
				}

				acked, err := r.clinet.XAck(ctx, r.stream, r.group, msg.ID).Result()
				if err != nil {
					logger.Log.Error(
						"Failed to acknowledge message",
						zap.String("message_id", msg.ID),
						zap.Error(err),
					)
					continue
				}

				logger.Log.Debug("Message acknowledged", zap.String("message_id", msg.ID), zap.Int64("acked", acked))
			}
		}
	}
}
