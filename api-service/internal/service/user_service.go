package service

import (
	"context"
	"errors"

	"github.com/excius/edns/internal/models"
	"github.com/excius/edns/internal/repository"
	"github.com/google/uuid"
)

var ErrUserExists = errors.New("service: user already exists")

type UserService struct {
	userRepo         *repository.UserRepository
	notificationRepo *repository.NotificationRepository
}

func NewUserService(userRepo *repository.UserRepository, notificationRepo *repository.NotificationRepository) *UserService {
	return &UserService{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserNotifications(ctx context.Context, userId uuid.UUID) ([]models.Notification, error) {
	return s.notificationRepo.ListUserNotifications(ctx, userId)
}

func (s *UserService) CreateUser(ctx context.Context, email string) (*models.User, error) {

	userExist, err := s.userRepo.UserExists(ctx, email)
	if err != nil {
		return nil, err
	}

	if userExist {
		return nil, ErrUserExists
	}

	user := &models.User{
		Email: email,
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
