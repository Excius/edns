package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationDelivery struct {
	ID             uuid.UUID  `json:"id"`
	NotificationID uuid.UUID  `json:"notification_id"`
	Channel        string     `json:"channel"`
	Status         string     `json:"status"`
	RetryCount     int        `json:"retry_count"`
	LastAttemptAt  *time.Time `json:"last_attempt_at"`
	CreatedAt      time.Time  `json:"created_at"`
}
