package processor

import (
	"context"

	"github.com/excius/edns/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationProcessor struct {
	db *pgxpool.Pool
}

func NewNotificationProcessor(db *pgxpool.Pool) *NotificationProcessor {
	return &NotificationProcessor{
		db: db,
	}
}

func (p *NotificationProcessor) Process(ctx context.Context, payload map[string]interface{}) {
	notificationID := payload["notification_id"].(string)

	logger.Log.Sugar().Infof("Processing notification with ID: %s", notificationID)

	// later we'll:
	// fetch delivery
	// send email
	// send websocket
	// update status
}
