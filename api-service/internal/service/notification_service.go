package service

import (
	"context"
	"errors"
	"strings"

	"github.com/excius/edns/pkg/models"
	"github.com/excius/edns/pkg/queue"
	"github.com/excius/edns/pkg/repository"
	"github.com/google/uuid"
)

type NotificationService struct {
	userRepo         *repository.UserRepository
	notificationRepo *repository.NotificationRepository
	deliveryRepo     *repository.NotificationDeliveryRepository
	queue            *queue.RedisStream
}

func NewNotificationService(
	userRepo *repository.UserRepository,
	notificationRepo *repository.NotificationRepository,
	deliveryRepo *repository.NotificationDeliveryRepository,
	queue *queue.RedisStream,
) *NotificationService {
	return &NotificationService{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		deliveryRepo:     deliveryRepo,
		queue:            queue,
	}
}

func (s *NotificationService) GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	return s.notificationRepo.GetNotificationByID(ctx, id)
}

func (s *NotificationService) GetDeliveriesByNotificatoinID(ctx context.Context, notificationID uuid.UUID) ([]models.NotificationDelivery, error) {
	return s.deliveryRepo.GetDeliveriesByNotificatoinID(ctx, notificationID)
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID uuid.UUID, message string, channels []string) (*models.Notification, error) {

	// Validate user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Validate channels
	var result []string
	seen := map[string]struct{}{}

	for _, ch := range channels {
		ch = strings.ToLower(strings.TrimSpace(ch))

		if !models.IsValidChannel(ch) {
			return nil, errors.New("invalid channel")
		}

		if _, exists := seen[ch]; !exists {
			seen[ch] = struct{}{}
			result = append(result, ch)
		}
	}

	// Create notification
	notification := &models.Notification{
		UserID:  userID,
		Message: message,
		Status:  models.StatusPending,
	}

	err = s.notificationRepo.CreateNotification(ctx, notification)
	if err != nil {
		return nil, err
	}

	// Create delivery records for each channel
	for _, channel := range result {

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

	// Publish to Redis Stream for processing
	err = s.queue.Publish(ctx, map[string]interface{}{
		"notification_id": notification.ID.String(),
	})

	if err != nil {
		return nil, err
	}

	return notification, nil
}
