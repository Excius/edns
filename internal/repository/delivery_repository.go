package repository

import (
	"context"

	"github.com/excius/edns/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationDeliveryRepository struct {
	db *pgxpool.Pool
}

func NewNotificationDeliveryRepository(db *pgxpool.Pool) *NotificationDeliveryRepository {
	return &NotificationDeliveryRepository{
		db: db,
	}
}

func (r *NotificationDeliveryRepository) CreateDelivery(ctx context.Context, delivery *models.NotificationDelivery) error {
	query := `
	INSERT INTO notification_deliveries (notification_id, channel, status)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query, delivery.NotificationID, delivery.Channel, delivery.Status).Scan(&delivery.ID, &delivery.CreatedAt)
}

func (r *NotificationDeliveryRepository) GetDeliveriesByNotificatoinID(ctx context.Context, id uuid.UUID) ([]models.NotificationDelivery, error) {
	query := `
	SELECT id, notification_id, channel, status, retry_count, last_attempt_at, created_at
	FROM notification_deliveries
	WHERE notification_id = $1
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []models.NotificationDelivery

	for rows.Next() {

		var d models.NotificationDelivery

		err := rows.Scan(&d.ID, &d.NotificationID, &d.Channel, &d.Status, &d.RetryCount, &d.LastAttemptAt, &d.CreatedAt)
		if err != nil {
			return nil, err
		}

		deliveries = append(deliveries, d)
	}

	return deliveries, nil
}

func (r *NotificationDeliveryRepository) GetByNotificationID(ctx context.Context, notificationID uuid.UUID) ([]models.NotificationDelivery, error) {
	query := `
	SELECT id, notification_id, channel, status, retry_count, last_attempt_at, created_at
	FROM notification_deliveries
	WHERE notification_id = $1
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, notificationID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []models.NotificationDelivery

	for rows.Next() {

		var d models.NotificationDelivery

		err := rows.Scan(&d.ID, &d.NotificationID, &d.Channel, &d.Status, &d.RetryCount, &d.LastAttemptAt, &d.CreatedAt)
		if err != nil {
			return nil, err
		}

		deliveries = append(deliveries, d)
	}

	return deliveries, nil
}

func (r *NotificationDeliveryRepository) UpdateStatus(ctx context.Context, deliveryID uuid.UUID, status string) error {
	query := `
	UPDATE notification_deliveries
	SET status = $1, last_attempt_at = NOW()
	WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, status, deliveryID)
	return err
}
