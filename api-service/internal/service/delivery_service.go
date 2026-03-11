package service

import (
	"context"
	"notification-system/api/internal/models"
	"notification-system/api/internal/repository"
)

type DeliveryService struct {
	deliveryRepo *repository.NotificationDeliveryRepository
}

func NewDeliveryService(delivertRepo *repository.NotificationDeliveryRepository) *DeliveryService {
	return &DeliveryService{
		deliveryRepo: delivertRepo,
	}
}

func (s *DeliveryService) CreateDelivery(ctx context.Context, delivery *models.NotificationDelivery) error {
	return s.deliveryRepo.CreateDelivery(ctx, delivery)
}
