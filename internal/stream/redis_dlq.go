package stream

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisDLQ struct {
	client *redis.Client
	stream string
}

func NewRedisDLQStream(client *redis.Client, stream string) *RedisDLQ {
	return &RedisDLQ{
		client: client,
		stream: stream,
	}
}

func (r *RedisDLQ) Publish(ctx context.Context, values map[string]any) error {

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.stream,
		Values: values,
	}).Err()
}
