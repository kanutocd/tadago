package repository

import (
	"github.com/google/uuid"
	"github.com/kanutocd/tada/internal/domain"
	"github.com/kanutocd/tada/internal/dto"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
	GetAll(pagination dto.PaginationQuery) ([]domain.User, string, error)
}

type TadaRepository interface {
	Create(tada *domain.Tada) error
	GetByID(id uuid.UUID) (*domain.Tada, error)
	Update(tada *domain.Tada) error
	Delete(id uuid.UUID) error
	GetAll(pagination dto.PaginationQuery) ([]domain.Tada, string, error)
	GetByUserID(userID uuid.UUID, pagination dto.PaginationQuery) ([]domain.Tada, string, error)
	GetByAssigneeID(assigneeID uuid.UUID, pagination dto.PaginationQuery) ([]domain.Tada, string, error)
}
