package queue

import "context"

type Processor interface {
	Process(ctx context.Context, payload map[string]any)
}
