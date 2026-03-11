package service

import (
	"context"
	"notification-system/api/internal/models"
	"notification-system/api/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email string) (*models.User, error) {

	user := &models.User{
		Email: email,
	}

	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
