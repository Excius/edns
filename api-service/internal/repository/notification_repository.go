package repository

import (
	"context"
	"notification-system/api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, notification *models.Notitication) error {
	query := `INSERT INTO notifications (user_id, message, status)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	return r.db.QueryRow(ctx, query, notification.UserID, notification.Message, notification.Status).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *NotificationRepository) GetNotification(ctx context.Context, id uuid.UUID) (*models.Notitication, error) {
	query := `SELCT id, user_id, message, status, created_at
	FROM notifications
	WHERE user_id = $1`

	row := r.db.QueryRow(ctx, query, id)

	var n models.Notitication

	err := row.Scan(&n.ID, &n.UserID, &n.Message, &n.Status, &n.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &n, nil
}
func (r *NotificationRepository) ListUserNotifications(ctx context.Context, userID uuid.UUID) ([]models.Notitication, error) {
	query := `SELECT id, user_id, message, status, created_at
	FROM notifications
	WHERE user_id = $1
	ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notitication

	for rows.Next() {
		var n models.Notitication
		err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.Status, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}
