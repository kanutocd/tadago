package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kanutocd/tada/internal/domain"
	"github.com/kanutocd/tada/internal/dto"
	"github.com/kanutocd/tada/internal/repository"
)

type UserService interface {
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(id uuid.UUID) (*dto.UserResponse, error)
	GetUsers(pagination dto.PaginationQuery) (*dto.PaginationResponse, error)
	UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	user := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return dto.ToUserResponse(user), nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return dto.ToUserResponse(user), nil
}

func (s *userService) GetUsers(pagination dto.PaginationQuery) (*dto.PaginationResponse, error) {
	users, nextCursor, err := s.userRepo.GetAll(pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *dto.ToUserResponse(&user)
	}

	return &dto.PaginationResponse{
		Data: userResponses,
		Pagination: dto.PaginationMeta{
			Limit:      pagination.Limit,
			Count:      len(userResponses),
			NextCursor: nextCursor,
		},
	}, nil
}

func (s *userService) UpdateUser(id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		// Check if new email already exists
		existingUser, err := s.userRepo.GetByEmail(*req.Email)
		if err == nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return dto.ToUserResponse(user), nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
