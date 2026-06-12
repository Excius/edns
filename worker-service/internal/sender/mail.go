package sender

import (
	"context"

	"github.com/excius/edns/internal/logger"
	"go.uber.org/zap"
)

type EmailSender struct{}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func (e *EmailSender) Send(ctx context.Context, notificationID string) error {
	logger.Log.Info(
		"Email sent",
		zap.String("notification_id", notificationID),
	)

	return nil
}
