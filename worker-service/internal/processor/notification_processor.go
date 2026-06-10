package processor

import (
	"context"
	"errors"
	"fmt"

	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/internal/models"
	"github.com/excius/edns/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const maxRetries = 3

type NotificationProcessor struct {
	notificationRepo *repository.NotificationRepository
	deliveryRepo     *repository.NotificationDeliveryRepository
}

func NewNotificationProcessor(
	notificationRepo *repository.NotificationRepository,
	deliveryRepo *repository.NotificationDeliveryRepository,
) *NotificationProcessor {
	return &NotificationProcessor{
		notificationRepo: notificationRepo,
		deliveryRepo:     deliveryRepo,
	}
}

func (p *NotificationProcessor) Process(
	ctx context.Context,
	payload map[string]any,
) error {

	notificationID, ok := payload["notification_id"].(string)
	if !ok {
		return errors.New("invalid payload: missing notification_id")
	}

	logger.Log.Info(
		"Processing notification",
		zap.String("notification_id", notificationID),
	)

	parsedNotiID, err := uuid.Parse(notificationID)
	if err != nil {
		return fmt.Errorf(
			"invalid notification id %s: %w",
			notificationID,
			err,
		)
	}

	deliveries, err := p.deliveryRepo.GetByNotificationID(ctx, parsedNotiID)
	if err != nil {
		return fmt.Errorf(
			"fetch deliveries for notification %s: %w",
			parsedNotiID,
			err,
		)
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
				return fmt.Errorf(
					"mark delivery %s failed: %w",
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
			return fmt.Errorf(
				"mark delivery %s processing: %w",
				delivery.ID,
				err,
			)
		}

		var deliveryErr error

		switch delivery.Channel {

		case "email":

			logger.Log.Info(
				"Processing email delivery",
				zap.String("notification_id", parsedNotiID.String()),
				zap.String("delivery_id", delivery.ID.String()),
			)

			// TODO: email sender
			deliveryErr = nil

		case "websocket":

			logger.Log.Info(
				"Processing websocket delivery",
				zap.String("notification_id", parsedNotiID.String()),
				zap.String("delivery_id", delivery.ID.String()),
			)

			// TODO: websocket sender
			deliveryErr = nil

		default:

			deliveryErr = fmt.Errorf(
				"unsupported delivery channel: %s",
				delivery.Channel,
			)
		}

		if deliveryErr != nil {

			retryCount, err := p.deliveryRepo.IncrementRetryCount(
				ctx,
				delivery.ID,
			)
			if err != nil {
				return fmt.Errorf(
					"increment retry count for delivery %s: %w",
					delivery.ID,
					err,
				)
			}

			logger.Log.Warn(
				"Delivery failed",
				zap.String("delivery_id", delivery.ID.String()),
				zap.Int("retry_count", retryCount),
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
				return fmt.Errorf(
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
			return fmt.Errorf(
				"mark delivery %s completed: %w",
				delivery.ID,
				err,
			)
		}
	}

	// TODO: Need to add Processing timeout recvery later in this
	if err := p.refreshNotificationStatus(
		ctx,
		parsedNotiID,
	); err != nil {
		return err
	}

	// For both failed and complete status we are pasing nil which will make ack for the message
	return nil
}

func (p *NotificationProcessor) refreshNotificationStatus(
	ctx context.Context,
	notificationID uuid.UUID,
) error {

	deliveries, err := p.deliveryRepo.GetByNotificationID(
		ctx,
		notificationID,
	)
	if err != nil {
		return fmt.Errorf(
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
		return fmt.Errorf(
			"update notification %s status: %w",
			notificationID,
			err,
		)
	}

	return nil
}
