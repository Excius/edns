package processor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/excius/edns/internal/events"
	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/internal/models"
	"github.com/excius/edns/internal/repository"
	"github.com/excius/edns/internal/stream"
	"github.com/excius/edns/worker-service/internal/metrics"
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
	dlqPublisher     *stream.RedisDLQ

	senders map[string]sender.Sender

	metrics *metrics.ProcessorMetrics
}

func NewNotificationProcessor(
	userRepo *repository.UserRepository,
	notificationRepo *repository.NotificationRepository,
	deliveryRepo *repository.NotificationDeliveryRepository,
	dlqPublisher *stream.RedisDLQ,
	senders map[string]sender.Sender,
	metrics *metrics.ProcessorMetrics,
) *NotificationProcessor {
	return &NotificationProcessor{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		deliveryRepo:     deliveryRepo,
		dlqPublisher:     dlqPublisher,
		senders:          senders,
		metrics:          metrics,
	}
}

// TODO:
// Prevent duplicate deliveries.
// Redis Streams provide at-least-once delivery, so retries may
// invoke the sender multiple times. Persist a delivery identifier
// (or external message ID) and skip sending if the delivery has
// already been completed.
func (p *NotificationProcessor) Process(
	ctx context.Context,
	payload map[string]any,
) (stream.ProcessResult, error) {

	start := time.Now()

	defer func() {
		p.metrics.NotificationProcessingDuration.Observe(
			time.Since(start).Seconds(),
		)
	}()

	p.metrics.NotificationsProcessed.Inc()

	notificationID, ok := payload["notification_id"].(string)
	if !ok {
		return stream.ProcessResult{}, errors.New("invalid payload: missing notification_id")
	}

	logger.FromContext(ctx).Info(
		"Processing notification",
		zap.String("notification_id", notificationID),
	)

	parsedNotiID, err := uuid.Parse(notificationID)
	if err != nil {
		return stream.ProcessResult{}, fmt.Errorf(
			"invalid notification id %s: %w",
			notificationID,
			err,
		)
	}

	deliveries, err := p.deliveryRepo.GetByNotificationID(ctx, parsedNotiID)
	if err != nil {
		return stream.ProcessResult{}, fmt.Errorf(
			"fetch deliveries for notification %s: %w",
			parsedNotiID,
			err,
		)
	}

	notificationContext, err := p.notificationRepo.GetNotificationWithUser(ctx, parsedNotiID)
	if err != nil {
		return stream.ProcessResult{}, fmt.Errorf(
			"fetch notification context for notification %s: %w",
			parsedNotiID,
			err,
		)
	}

	event := events.NotificationEvent{
		NotificationID: parsedNotiID.String(),
		UserID:         notificationContext.User.ID.String(),
		Email:          notificationContext.User.Email,
		Title:          notificationContext.Notfication.Title,
		Message:        notificationContext.Notfication.Message,
	}

	for _, delivery := range deliveries {

		// Already finished
		if delivery.Status == models.StatusCompleted ||
			delivery.Status == models.StatusFailed {
			continue
		}

		p.metrics.DeliveriesProcessed.Inc()

		// Retry budget exhausted
		if delivery.RetryCount >= maxRetries {

			if err := p.deliveryRepo.UpdateStatus(
				ctx,
				delivery.ID,
				models.StatusFailed,
			); err != nil {
				return stream.ProcessResult{}, fmt.Errorf(
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
				return stream.ProcessResult{}, fmt.Errorf(
					"publish delivery %s to dlq: %w",
					delivery.ID,
					err,
				)
			}

			p.metrics.DLQMessages.Inc()

			continue
		}

		if err := p.deliveryRepo.UpdateStatus(
			ctx,
			delivery.ID,
			models.StatusProcessing,
		); err != nil {
			return stream.ProcessResult{}, fmt.Errorf(
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

			deliveryStart := time.Now()

			deliveryErr = deliverySender.Send(ctx, event)

			p.metrics.DeliveryDuration.Observe(
				time.Since(deliveryStart).Seconds(),
			)
		}

		if deliveryErr != nil {

			retryCount, err := p.deliveryRepo.IncrementRetryCount(
				ctx,
				delivery.ID,
			)
			if err != nil {
				return stream.ProcessResult{}, fmt.Errorf(
					"increment retry count for delivery %s: %w",
					delivery.ID,
					err,
				)
			}

			p.metrics.DeliveryRetries.Inc()

			logger.FromContext(ctx).Warn(
				"Delivery failed",
				zap.String("delivery_id", delivery.ID.String()),
				zap.Int("retry_count", retryCount),
				zap.Error(deliveryErr),
			)

			status := models.StatusPending
			if retryCount >= maxRetries {
				p.metrics.DeliveriesFailed.Inc()
				status = models.StatusFailed
			}

			if err := p.deliveryRepo.UpdateStatus(
				ctx,
				delivery.ID,
				status,
			); err != nil {
				return stream.ProcessResult{}, fmt.Errorf(
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
			return stream.ProcessResult{}, fmt.Errorf(
				"mark delivery %s completed: %w",
				delivery.ID,
				err,
			)
		}

		p.metrics.DeliveriesCompleted.Inc()
	}

	// TODO:
	// Recover deliveries stuck in StatusProcessing for longer than the configured timeout.
	// Periodically scan the database, mark stale deliveries back to StatusPending,
	// increment the retry count, and let Redis Stream recovery retry them.
	status, err := p.refreshNotificationStatus(
		ctx,
		parsedNotiID,
	)
	if err != nil {
		return stream.ProcessResult{}, err
	}

	// Sending the reuqest to retry state
	if status != models.StatusCompleted && status != models.StatusFailed {
		return stream.ProcessResult{
			Ack: false,
		}, nil
	}

	// For both failed and complete status we are pasing nil which will make ack for the message

	if status == models.StatusCompleted {
		p.metrics.NotificationsCompleted.Inc()
	} else {
		p.metrics.NotificationsFailed.Inc()
	}

	return stream.ProcessResult{
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
