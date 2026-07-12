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
	messageCount     int64
	recoveryInterval int
	recoveryIdleTime int
}

func NewRedisRecovery(
	client *redis.Client,
	stream string,
	group string,
	consumerName string,
	messageCount int64,
	recoveryInterval int,
	recoveryIdleTime int,
) *RedisRecovery {
	return &RedisRecovery{
		client:           client,
		stream:           stream,
		group:            group,
		consumerName:     consumerName,
		messageCount:     messageCount,
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

				select {
				case <-ctx.Done():
					return nil

				default:
				}

				msgs, nextCursor, err := r.client.XAutoClaim(
					ctx,
					&redis.XAutoClaimArgs{
						Stream:   r.stream,
						Group:    r.group,
						MinIdle:  time.Duration(r.recoveryIdleTime) * time.Second,
						Consumer: r.consumerName,
						Count:    r.messageCount,
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

					// Before processing each message
					select {
					case <-ctx.Done():
						return nil
					default:
					}

					logger.Log.Info(
						"Recovered stale message",
						zap.String("message_id", msg.ID),
					)

					result, err := processor.Process(ctx, msg.Values)
					if err != nil {
						logger.Log.Error(
							"Failed processing message",
							zap.String("message_id", msg.ID),
							zap.Error(err),
						)

						// Leave message in PEL
						continue
					}

					// Leaving for recovery to retry this message
					if !result.Ack {
						logger.Log.Debug(
							"Notification still in progress",
							zap.String("message_id", msg.ID),
						)

						// Leave message in PEL
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
						zap.Int64("acked", acked),
					)
				}

				// Break out of the infinite for loop
				if nextCursor == "0-0" {
					break recoveryLoop
				}

				cursor = nextCursor
			}
		}
	}
}
