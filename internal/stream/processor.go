package stream

import "context"

type ProcessResult struct {
	// Ack indicates whether the Redis Stream message
	// can be acknowledged.
	Ack bool
}

type Processor interface {
	Process(ctx context.Context, payload map[string]any) (ProcessResult, error)
}
