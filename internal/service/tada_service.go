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

type TadaService interface {
	CreateTada(req dto.CreateTadaRequest) (*dto.TadaResponse, error)
	GetTadaByID(id uuid.UUID) (*dto.TadaResponse, error)
	GetTadas(pagination dto.PaginationQuery) (*dto.PaginationResponse, error)
	UpdateTada(id uuid.UUID, req dto.UpdateTadaRequest) (*dto.TadaResponse, error)
	DeleteTada(id uuid.UUID) error
}

type tadaService struct {
	tadaRepo repository.TadaRepository
	userRepo repository.UserRepository
}

func NewTadaService(tadaRepo repository.TadaRepository, userRepo repository.UserRepository) TadaService {
	return &tadaService{
		tadaRepo: tadaRepo,
		userRepo: userRepo,
	}
}

func (s *tadaService) CreateTada(req dto.CreateTadaRequest) (*dto.TadaResponse, error) {
	// Validate creator exists
	_, err := s.userRepo.GetByID(req.CreatedBy)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("creator user not found")
		}
		return nil, fmt.Errorf("failed to validate creator: %w", err)
	}

	// Validate assignee exists if provided
	if req.AssignedTo != nil {
		_, err := s.userRepo.GetByID(*req.AssignedTo)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("assignee user not found")
			}
			return nil, fmt.Errorf("failed to validate assignee: %w", err)
		}
	}

	tada := &domain.Tada{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
		AssignedTo:  req.AssignedTo,
		DueAt:       req.DueAt,
	}

	if req.Status != nil {
		tada.Status = *req.Status
	}

	if err := s.tadaRepo.Create(tada); err != nil {
		return nil, fmt.Errorf("failed to create tada: %w", err)
	}

	// Reload with relationships
	tada, _ = s.tadaRepo.GetByID(tada.ID)

	return dto.ToTadaResponse(tada), nil
}

func (s *tadaService) GetTadaByID(id uuid.UUID) (*dto.TadaResponse, error) {
	tada, err := s.tadaRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tada not found")
		}
		return nil, fmt.Errorf("failed to get tada: %w", err)
	}

	return dto.ToTadaResponse(tada), nil
}

func (s *tadaService) GetTadas(pagination dto.PaginationQuery) (*dto.PaginationResponse, error) {
	tadas, nextCursor, err := s.tadaRepo.GetAll(pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get tadas: %w", err)
	}

	tadaResponses := make([]dto.TadaResponse, len(tadas))
	for i, tada := range tadas {
		tadaResponses[i] = *dto.ToTadaResponse(&tada)
	}

	return &dto.PaginationResponse{
		Data: tadaResponses,
		Pagination: dto.PaginationMeta{
			Limit:      pagination.Limit,
			Count:      len(tadaResponses),
			NextCursor: nextCursor,
		},
	}, nil
}

func (s *tadaService) UpdateTada(id uuid.UUID, req dto.UpdateTadaRequest) (*dto.TadaResponse, error) {
	tada, err := s.tadaRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tada not found")
		}
		return nil, fmt.Errorf("failed to get tada: %w", err)
	}

	// Update fields
	if req.Name != nil {
		tada.Name = *req.Name
	}
	if req.Description != nil {
		tada.Description = *req.Description
	}
	if req.AssignedTo != nil {
		// Validate assignee exists
		_, err := s.userRepo.GetByID(*req.AssignedTo)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("assignee user not found")
			}
			return nil, fmt.Errorf("failed to validate assignee: %w", err)
		}
		tada.AssignedTo = req.AssignedTo
	}
	if req.Status != nil {
		tada.Status = *req.Status
	}
	if req.DueAt != nil {
		tada.DueAt = req.DueAt
	}

	if err := s.tadaRepo.Update(tada); err != nil {
		return nil, fmt.Errorf("failed to update tada: %w", err)
	}

	// Reload with relationships
	tada, _ = s.tadaRepo.GetByID(tada.ID)

	return dto.ToTadaResponse(tada), nil
}

func (s *tadaService) DeleteTada(id uuid.UUID) error {
	_, err := s.tadaRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tada not found")
		}
		return fmt.Errorf("failed to get tada: %w", err)
	}

	if err := s.tadaRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete tada: %w", err)
	}

	return nil
}
