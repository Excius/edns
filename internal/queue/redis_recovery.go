package queue

import (
	"context"
	"time"

	"github.com/excius/edns/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisRecovery struct {
	client           *redis.Client
	stream           string
	group            string
	consumerName     string
	recoveryInterval int
	recoveryIdleTime int
}

func NewRedisRecovery(client *redis.Client, stream string, group string, consumerName string, recoveryInterval int, recoveryIdleTime int) *RedisRecovery {
	return &RedisRecovery{
		client:           client,
		stream:           stream,
		group:            group,
		consumerName:     consumerName,
		recoveryInterval: recoveryInterval, // after what duration recovery runs
		recoveryIdleTime: recoveryIdleTime, // how old msg needs to be processed
	}
}

func (r *RedisRecovery) Start(ctx context.Context, processor Processor) error {

	// timeout in which the recovery starts
	ticker := time.NewTicker(time.Duration(r.recoveryInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:

			cursor := "0-0"

		recoveryLoop:
			for {
				// TODO: Need to make the count a config var
				msgs, nextCursor, err := r.client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
					Stream:   r.stream,
					Group:    r.group,
					MinIdle:  time.Duration(r.recoveryIdleTime) * time.Second,
					Consumer: r.consumerName,
					Count:    4,
					Start:    cursor,
				}).Result()

				if err != nil {
					logger.Log.Error(
						"XAUTOCLAIM failed",
						zap.Error(err),
					)
					break recoveryLoop
				}

				logger.Log.Debug(
					"Recovery scan",
					zap.String("cursor", cursor),
					zap.String("next_cursor", nextCursor),
					zap.Int("messages", len(msgs)),
				)

				if len(msgs) == 0 {
					break recoveryLoop
				}

				for _, msg := range msgs {

					logger.Log.Info(
						"Recovered stale message",
						zap.String("message_id", msg.ID),
					)

					if err := processor.Process(ctx, msg.Values); err != nil {

						logger.Log.Error(
							"Failed processing message",
							zap.String("message_id", msg.ID),
							zap.Error(err),
						)
						continue
					}

					acked, err := r.client.XAck(ctx, r.stream, r.group, msg.ID).Result()
					if err != nil {
						logger.Log.Error(
							"Failed to acknowledge message",
							zap.String("message_id", msg.ID),
							zap.Error(err),
						)
						continue
					}

					logger.Log.Debug(
						"Message acknowledged",
						zap.String("message_id", msg.ID),
						zap.Int64("acked", acked))
				}

				if nextCursor == "0-0" {
					break recoveryLoop
				}

				cursor = nextCursor
			}
		}
	}
}
