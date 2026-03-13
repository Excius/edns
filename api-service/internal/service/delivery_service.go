package service

import (
	"context"

	"github.com/excius/edns/pkg/models"
	"github.com/excius/edns/pkg/repository"
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
