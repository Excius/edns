package sender

import (
	"context"

	"github.com/excius/edns/internal/events"
)

type Sender interface {
	Send(ctx context.Context, event events.NotificationEvent) error
}
