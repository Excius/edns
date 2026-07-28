package service

import (
	"context"
	"errors"

	"github.com/excius/edns/api-service/internal/metrics"
	"github.com/excius/edns/internal/apperrors"
	"github.com/excius/edns/internal/models"
	"github.com/excius/edns/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	userRepo         *repository.UserRepository
	notificationRepo *repository.NotificationRepository
	metrics          *metrics.Metrics
}

func NewUserService(userRepo *repository.UserRepository, notificationRepo *repository.NotificationRepository, metrics *metrics.Metrics) *UserService {
	return &UserService{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		metrics:          metrics,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserNotifications(ctx context.Context, userId uuid.UUID) ([]models.Notification, error) {
	notifications, err := s.notificationRepo.ListUserNotifications(ctx, userId)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *UserService) CreateUser(ctx context.Context, email string) (*models.User, error) {

	userExist, err := s.userRepo.UserExists(ctx, email)
	if err != nil {
		return nil, err
	}

	if userExist {
		return nil, apperrors.ErrUserExists
	}

	user := &models.User{
		Email: email,
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	s.metrics.Business.UsersCreated.Inc()

	return user, nil
}
