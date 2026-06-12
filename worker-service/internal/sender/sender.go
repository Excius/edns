package sender

import (
	"context"
)

type Sender interface {
	Send(ctx context.Context, notificationID string) error
}
