package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisConsumer struct {
	clinet *redis.Client
	stream string
}

func NewRedisConsumer(client *redis.Client, stream string) *RedisConsumer {
	return &RedisConsumer{
		clinet: client,
		stream: stream,
	}
}

func (r *RedisConsumer) Start(ctx context.Context, processor Processor) error {

	for {

		streams, err := r.clinet.XRead(ctx, &redis.XReadArgs{
			Streams: []string{r.stream, "0"},
			Block:   0,
			Count:   1,
		}).Result()

		if err != nil {
			return err
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {

				processor.Process(ctx, msg.Values)

			}
		}
	}
}
