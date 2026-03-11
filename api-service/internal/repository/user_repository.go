package repository

import (
	"context"
	"notification-system/api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
	INSERT INTO users (email)
	VALUES ($1)
	RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query, user.Email).Scan(&user.ID, &user.CretedAt)
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
	SELECT id, email, created_at
	FROM users
	WHERE id = $1
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.CretedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
