package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kanutocd/tada/internal/domain"
)

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=255"`
	Email string `json:"email" binding:"required,email,max=255"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Email *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
