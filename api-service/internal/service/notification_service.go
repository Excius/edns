package service

import (
	"context"
	"notification-system/api/internal/models"
	"notification-system/api/internal/repository"

	"github.com/google/uuid"
)

type NotificationService struct {
	userRepo         *repository.UserRepository
	notificationRepo *repository.NotificationRepository
	deliveryRepo     *repository.NotificationDeliveryRepository
}

func NewNotificationService(
	userRepo *repository.UserRepository,
	notificationRepo *repository.NotificationRepository,
	deliveryRepo *repository.NotificationDeliveryRepository,
) *NotificationService {
	return &NotificationService{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		deliveryRepo:     deliveryRepo,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID uuid.UUID, message string, channels []string) (*models.Notitication, error) {

	// Validate user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create notification
	notification := &models.Notitication{
		UserID:  userID,
		Message: message,
		Status:  models.StatusPending,
	}

	err = s.notificationRepo.CreateNotification(ctx, notification)
	if err != nil {
		return nil, err
	}

	// Create delivery records for each channel
	for _, channel := range channels {

		delivery := &models.NotificationDelivery{
			NotificationID: notification.ID,
			Channel:        channel,
			Status:         models.StatusPending,
		}

		err := s.deliveryRepo.CreateDelivery(ctx, delivery)
		if err != nil {
			return nil, err
		}

	}

	return notification, nil
}
