package subscriber

import "context"

type EventHandler interface {
	Handle(ctx context.Context, payload []byte) error
}
