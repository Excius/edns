package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStream struct {
	client *redis.Client
	stream string
}

func NewRedisStream(client *redis.Client, stream string) *RedisStream {
	return &RedisStream{
		client: client,
		stream: stream,
	}
}

func (r *RedisStream) Publish(ctx context.Context, values map[string]any) error {

	return r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.stream,
		Values: values,
	}).Err()
}
