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

func (r *UserRepository) UserExists(ctx context.Context, email string) (bool, error) {
	query := `
	SELECT COUNT(*)
	FROM users
	WHERE email = $1
	`

	var count int

	err := r.db.QueryRow(ctx, query, email).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
	SELECT id, email, created_at
	FROM users
	WHERE email = $1
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.CretedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
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
