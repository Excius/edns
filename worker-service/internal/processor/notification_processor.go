package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/excius/edns/internal/events"
	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/internal/models"
	"github.com/excius/edns/internal/queue"
	"github.com/excius/edns/internal/repository"
	"github.com/excius/edns/worker-service/internal/sender"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ErrNotificationPending = errors.New("notification pending")

const maxRetries = 3

type NotificationProcessor struct {
	userRepo         *repository.UserRepository
	notificationRepo *repository.NotificationRepository
	deliveryRepo     *repository.NotificationDeliveryRepository
	dlqPublisher     *queue.RedisDLQ

	senders map[string]sender.Sender
}

func NewNotificationProcessor(
	userRepo *repository.UserRepository,
	notificationRepo *repository.NotificationRepository,
	deliveryRepo *repository.NotificationDeliveryRepository,
	dlqPublisher *queue.RedisDLQ,
	senders map[string]sender.Sender,
) *NotificationProcessor {
	return &NotificationProcessor{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		deliveryRepo:     deliveryRepo,
		dlqPublisher:     dlqPublisher,
		senders:          senders,
	}
}

// TODO: Add idempotency protection when real email/websocket
// delivery is implemented. Redis Streams provide at-least-once delivery.
func (p *NotificationProcessor) Process(
	ctx context.Context,
	payload map[string]any,
) (queue.ProcessResult, error) {

	notificationID, ok := payload["notification_id"].(string)
	if !ok {
		return queue.ProcessResult{}, errors.New("invalid payload: missing notification_id")
	}

	logger.Log.Info(
		"Processing notification",
		zap.String("notification_id", notificationID),
	)

	parsedNotiID, err := uuid.Parse(notificationID)
	if err != nil {
		return queue.ProcessResult{}, fmt.Errorf(
			"invalid notification id %s: %w",
			notificationID,
			err,
		)
	}

	deliveries, err := p.deliveryRepo.GetByNotificationID(ctx, parsedNotiID)
	if err != nil {
		return queue.ProcessResult{}, fmt.Errorf(
			"fetch deliveries for notification %s: %w",
			parsedNotiID,
			err,
		)
	}

	// TODO: can make a join request to get data in one request rather then 3 simulateneous requests

	notification, err := p.notificationRepo.GetNotificationByID(ctx, parsedNotiID)
	if err != nil {
		return queue.ProcessResult{}, fmt.Errorf(
			"fetch notification for notification %s: %w",
			parsedNotiID,
			err,
		)
	}

	user, err := p.userRepo.GetUserByID(ctx, notification.UserID)
	if err != nil {
		return queue.ProcessResult{}, fmt.Errorf(
			"fetch user details for notification %s: %w",
			parsedNotiID,
			err,
		)
	}

	event := events.NotificationEvent{
		NotificationID: parsedNotiID.String(),
		UserID:         notification.UserID.String(),
		Email:          user.Email,
		Title:          notification.Title,
		Message:        notification.Message,
	}

	for _, delivery := range deliveries {

		// Already finished
		if delivery.Status == models.StatusCompleted ||
			delivery.Status == models.StatusFailed {
			continue
		}

		// Retry budget exhausted
		if delivery.RetryCount >= maxRetries {

			if err := p.deliveryRepo.UpdateStatus(
				ctx,
				delivery.ID,
				models.StatusFailed,
			); err != nil {
				return queue.ProcessResult{}, fmt.Errorf(
					"mark delivery %s failed: %w",
					delivery.ID,
					err,
				)
			}

			// Add to dlq stream
			err := p.dlqPublisher.Publish(ctx, map[string]any{
				"notification_id": parsedNotiID.String(),
				"delivery_id":     delivery.ID.String(),
				"channel":         delivery.Channel,
				"reason":          "max retries exceeded",
			})

			if err != nil {
				return queue.ProcessResult{}, fmt.Errorf(
					"publish delivery %s to dlq: %w",
					delivery.ID,
					err,
				)
			}

			continue
		}

		if err := p.deliveryRepo.UpdateStatus(
			ctx,
			delivery.ID,
			models.StatusProcessing,
		); err != nil {
			return queue.ProcessResult{}, fmt.Errorf(
				"mark delivery %s processing: %w",
				delivery.ID,
				err,
			)
		}

		var deliveryErr error

		deliverySender, exists := p.senders[delivery.Channel]
		if !exists {
			deliveryErr = fmt.Errorf(
				"unsupported delivery channel: %s",
				delivery.Channel,
			)
		} else {
			deliveryErr = deliverySender.Send(ctx, event)
		}

		if deliveryErr != nil {

			retryCount, err := p.deliveryRepo.IncrementRetryCount(
				ctx,
				delivery.ID,
			)
			if err != nil {
				return queue.ProcessResult{}, fmt.Errorf(
					"increment retry count for delivery %s: %w",
					delivery.ID,
					err,
				)
			}

			logger.Log.Warn(
				"Delivery failed",
				zap.String("delivery_id", delivery.ID.String()),
				zap.Int("retry_count", retryCount),
				zap.Error(deliveryErr),
			)

			status := models.StatusPending
			if retryCount >= maxRetries {
				status = models.StatusFailed
			}

			if err := p.deliveryRepo.UpdateStatus(
				ctx,
				delivery.ID,
				status,
			); err != nil {
				return queue.ProcessResult{}, fmt.Errorf(
					"mark delivery %s pending: %w",
					delivery.ID,
					err,
				)
			}

			continue
		}

		if err := p.deliveryRepo.UpdateStatus(
			ctx,
			delivery.ID,
			models.StatusCompleted,
		); err != nil {
			return queue.ProcessResult{}, fmt.Errorf(
				"mark delivery %s completed: %w",
				delivery.ID,
				err,
			)
		}
	}

	// TODO: Need to add Processing timeout recvery later in this
	status, err := p.refreshNotificationStatus(
		ctx,
		parsedNotiID,
	)
	if err != nil {
		return queue.ProcessResult{}, err
	}

	if status != models.StatusCompleted && status != models.StatusFailed {
		return queue.ProcessResult{
			Ack: false,
		}, nil
	}

	// For both failed and complete status we are pasing nil which will make ack for the message
	return queue.ProcessResult{
		Ack: true,
	}, nil
}

func (p *NotificationProcessor) refreshNotificationStatus(
	ctx context.Context,
	notificationID uuid.UUID,
) (string, error) {
	deliveries, err := p.deliveryRepo.GetByNotificationID(
		ctx,
		notificationID,
	)
	if err != nil {
		return models.StatusPending, fmt.Errorf(
			"fetch deliveries for notification status refresh %s: %w",
			notificationID,
			err,
		)
	}

	hasFailed := false
	hasPending := false

	for _, delivery := range deliveries {
		switch delivery.Status {

		case models.StatusPending, models.StatusProcessing:
			hasPending = true

		case models.StatusFailed:
			hasFailed = true
		}
	}

	status := models.StatusCompleted

	switch {
	case hasPending:
		status = models.StatusPending

	case hasFailed:
		status = models.StatusFailed
	}

	if err := p.notificationRepo.UpdateNotificationStatus(
		ctx,
		notificationID,
		status,
	); err != nil {
		return models.StatusPending, fmt.Errorf(
			"update notification %s status: %w",
			notificationID,
			err,
		)
	}

	return status, nil
}
