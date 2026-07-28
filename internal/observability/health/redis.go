package health

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisChecker struct {
	client *redis.Client
}

func NewRedisChecker(client *redis.Client) *RedisChecker {
	return &RedisChecker{
		client: client,
	}
}

func (r *RedisChecker) Check(ctx context.Context) CheckResult {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return CheckResult{
			Name:   "redis",
			Status: "down",
			Error:  err.Error(),
		}
	}

	return CheckResult{
		Name:   "redis",
		Status: "up",
	}
}
