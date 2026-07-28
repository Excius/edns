package repository

import (
	"context"

	"github.com/excius/edns/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationContext struct {
	Notfication models.Notification
	User        models.User
}

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	query := `INSERT INTO notifications (user_id, title, message, status)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query, notification.UserID, notification.Title, notification.Message, notification.Status).Scan(&notification.ID, &notification.CreatedAt)
}

func (r *NotificationRepository) GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	query := `SELECT id, user_id, title, message, status, created_at
	FROM notifications
	WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var n models.Notification

	err := row.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Status, &n.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *NotificationRepository) GetNotificationByUserID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	query := `SELECT id, user_id, title, message, status, created_at
	FROM notifications
	WHERE user_id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var n models.Notification

	err := row.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Status, &n.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (r *NotificationRepository) ListUserNotifications(ctx context.Context, userID uuid.UUID) ([]models.Notification, error) {
	query := `SELECT id, user_id, title, message, status, created_at
	FROM notifications
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var n models.Notification
		err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Status, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *NotificationRepository) UpdateNotificationStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE notifications
	SET status = $1
	WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *NotificationRepository) GetNotificationWithUser(ctx context.Context, notificationID uuid.UUID) (*NotificationContext, error) {
	query := `SELECT
	n.id,
	n.title,
	n.message,
	n.status,
	n.created_at,

	u.id,
	u.email,
	u.created_at
	FROM notifications n
	JOIN users u
	ON u.id = n.user_id
	WHERE n.id = $1;
	`

	row := r.db.QueryRow(ctx, query, notificationID)

	var n NotificationContext

	err := row.Scan(
		&n.Notfication.ID,
		&n.Notfication.Title,
		&n.Notfication.Message,
		&n.Notfication.Status,
		&n.Notfication.CreatedAt,
		&n.User.ID,
		&n.User.Email,
		&n.User.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &n, nil

}
