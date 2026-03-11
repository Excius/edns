package dto

import "github.com/google/uuid"

type CreateNotificationRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	Message  string    `json:"message" binding:"required"`
	Channels []string  `json:"channels" binding:"required"`
}
