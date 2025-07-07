package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/kanutocd/tada/internal/domain"
)

type CreateTadaRequest struct {
	Name        string             `json:"name" binding:"required,min=1,max=255"`
	Description string             `json:"description"`
	CreatedBy   uuid.UUID          `json:"created_by" binding:"required"`
	AssignedTo  *uuid.UUID         `json:"assigned_to,omitempty"`
	Status      *domain.TadaStatus `json:"status,omitempty" binding:"omitempty,oneof=in_progress cancelled completed"`
	DueAt       *time.Time         `json:"due_at,omitempty"`
}

type UpdateTadaRequest struct {
	Name        *string            `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string            `json:"description,omitempty"`
	AssignedTo  *uuid.UUID         `json:"assigned_to,omitempty"`
	Status      *domain.TadaStatus `json:"status,omitempty" binding:"omitempty,oneof=in_progress cancelled completed"`
	DueAt       *time.Time         `json:"due_at,omitempty"`
}

type TadaResponse struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CreatedBy   uuid.UUID         `json:"created_by"`
	AssignedTo  *uuid.UUID        `json:"assigned_to,omitempty"`
	Status      domain.TadaStatus `json:"status"`
	DueAt       *time.Time        `json:"due_at,omitempty"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	Creator     *UserResponse     `json:"creator,omitempty"`
	Assignee    *UserResponse     `json:"assignee,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ToTadaResponse(tada *domain.Tada) *TadaResponse {
	response := &TadaResponse{
		ID:          tada.ID,
		Name:        tada.Name,
		Description: tada.Description,
		CreatedBy:   tada.CreatedBy,
		AssignedTo:  tada.AssignedTo,
		Status:      tada.Status,
		DueAt:       tada.DueAt,
		CompletedAt: tada.CompletedAt,
		CreatedAt:   tada.CreatedAt,
		UpdatedAt:   tada.UpdatedAt,
	}

	if tada.Creator.ID != uuid.Nil {
		response.Creator = ToUserResponse(&tada.Creator)
	}

	if tada.Assignee != nil && tada.Assignee.ID != uuid.Nil {
		response.Assignee = ToUserResponse(tada.Assignee)
	}

	return response
}
