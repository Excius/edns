package events

type NotificationEvent struct {
	NotificationID string `json:"notification_id"`
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	Title          string `json:"title"`
	Message        string `json:"message"`
}
