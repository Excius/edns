package dto

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}
