package processor

import (
	"context"
	"database/sql"
	"errors"

	"github.com/excius/edns/internal/logger"
	"github.com/excius/edns/internal/models"
	"github.com/excius/edns/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationProcessor struct {
	notificationRepo *repository.NotificationRepository
	deliveryRepo     *repository.NotificationDeliveryRepository
}

func NewNotificationProcessor(notificationRepo *repository.NotificationRepository, deliveryRepo *repository.NotificationDeliveryRepository) *NotificationProcessor {
	return &NotificationProcessor{
		notificationRepo: notificationRepo,
		deliveryRepo:     deliveryRepo,
	}
}

func (p *NotificationProcessor) Process(ctx context.Context, payload map[string]interface{}) {
	notificationID, ok := payload["notification_id"].(string)
	if !ok {
		logger.Log.Error("Invalid payload: missing notification_id")
		return
	}

	logger.Log.Info("Processing notification", zap.String("notification_id", notificationID))

	parsedNotiID, err := uuid.Parse(notificationID)
	logger.Log.Error("Invalid notification ID format", zap.String("notification_id", notificationID), zap.Error(err))

	_, err = p.notificationRepo.GetNotificationByID(ctx, parsedNotiID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Warn("Notification not found", zap.String("notification_id", parsedNotiID.String()))
		} else {
			logger.Log.Error("Failed to fetch notification", zap.String("notification_id", parsedNotiID.String()), zap.Error(err))
		}
		return
	}

	deliveries, err := p.deliveryRepo.GetByNotificationID(ctx, parsedNotiID)
	if err != nil {
		logger.Log.Error("Failed to fetch deliveries", zap.String("notification_id", parsedNotiID.String()), zap.Error(err))
		return
	}

	allDelivered := true

	// Process each delivery
	for _, delivery := range deliveries {

		if delivery.Status == models.StatusCompleted {
			continue
		}

		// Mark as processing
		err := p.deliveryRepo.UpdateStatus(ctx, delivery.ID, models.StatusProcessing)
		if err != nil {
			logger.Log.Error("Failed to mark delivery processing", zap.String("delivery_id", delivery.ID.String()), zap.Error(err))
			allDelivered = false
			continue
		}

		var deliveryErr error

		switch delivery.Channel {

		case "email":
			logger.Log.Info("Processing email delivery", zap.String("notification_id", parsedNotiID.String()), zap.String("delivery_id", delivery.ID.String()))
			deliveryErr = nil // TODO: implement email sending logic

		case "websocket":
			logger.Log.Info("Processing websocket delivery", zap.String("notification_id", parsedNotiID.String()), zap.String("delivery_id", delivery.ID.String()))
			deliveryErr = nil // TODO: implement websocket sending logic

		default:
			logger.Log.Warn("Unknown delivery channel", zap.String("channel", delivery.Channel), zap.String("notification_id", parsedNotiID.String()))
			deliveryErr = errors.New("unknown delivery channel")
		}

		// Update final delivery status
		status := models.StatusCompleted
		if deliveryErr != nil {
			allDelivered = false
			status = models.StatusFailed
		}

		err = p.deliveryRepo.UpdateStatus(ctx, delivery.ID, status)
		if err != nil {
			logger.Log.Error("Failed to update delivery status", zap.String("delivery_id", delivery.ID.String()), zap.String("status", status), zap.Error(err))
		}
	}

	finalStatus := models.StatusCompleted
	if !allDelivered {
		finalStatus = models.StatusFailed
	}

	err = p.notificationRepo.UpdateNotificationStatus(ctx, parsedNotiID, finalStatus)
	if err != nil {
		logger.Log.Error("Failed to update notification status", zap.String("notification_id", parsedNotiID.String()), zap.String("status", finalStatus), zap.Error(err))
	}
}
