package apperrors

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrNotificationNotFound = errors.New("notification not found")
	ErrDeliveryNotFound     = errors.New("delivery not found")

	ErrInvalidChannel = errors.New("invalid notification channel")

	ErrUserExists = errors.New("user already exists")
)
