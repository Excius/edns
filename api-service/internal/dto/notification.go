package dto

type CreateNotificationRequest struct {
	UserID   string   `json:"user_id" binding:"required"`
	Title    string   `json:"title" binding:"required,max=80"`
	Message  string   `json:"message" binding:"required,max=500"`
	Channels []string `json:"channels" binding:"required"`
}
