package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kanutocd/tada/internal/domain"
	"github.com/kanutocd/tada/internal/dto"
)

type tadaRepository struct {
	db *gorm.DB
}

func NewTadaRepository(db *gorm.DB) TadaRepository {
	return &tadaRepository{db: db}
}

func (r *tadaRepository) Create(tada *domain.Tada) error {
	return r.db.Create(tada).Error
}

func (r *tadaRepository) GetByID(id uuid.UUID) (*domain.Tada, error) {
	var tada domain.Tada
	err := r.db.Preload("Creator").Preload("Assignee").First(&tada, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tada, nil
}

func (r *tadaRepository) Update(tada *domain.Tada) error {
	return r.db.Save(tada).Error
}

func (r *tadaRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Tada{}, "id = ?", id).Error
}

func (r *tadaRepository) GetAll(pagination dto.PaginationQuery) ([]domain.Tada, string, error) {
	var tadas []domain.Tada
	var query *gorm.DB = r.db.Model(&domain.Tada{}).Preload("Creator").Preload("Assignee")

	// Set default limit
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	// Apply cursor pagination
	if pagination.Cursor != "" {
		cursor, err := dto.DecodeCursor(pagination.Cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor: %w", err)
		}

		query = query.Where("(created_at < ? OR (created_at = ? AND id < ?))",
			cursor.CreatedAt, cursor.CreatedAt, cursor.ID)
	}

	// Fetch one extra record to check if there are more pages
	query = query.Order("created_at DESC, id DESC").Limit(pagination.Limit + 1)

	if err := query.Find(&tadas).Error; err != nil {
		return nil, "", fmt.Errorf("failed to fetch tadas: %w", err)
	}

	// Generate next cursor
	var nextCursor string
	if len(tadas) > pagination.Limit {
		tadas = tadas[:pagination.Limit] // Remove extra record
		lastTada := tadas[len(tadas)-1]
		nextCursor = dto.EncodeCursor(dto.Cursor{
			ID:        lastTada.ID.String(),
			CreatedAt: lastTada.CreatedAt,
		})
	}

	return tadas, nextCursor, nil
}

func (r *tadaRepository) GetByUserID(userID uuid.UUID, pagination dto.PaginationQuery) ([]domain.Tada, string, error) {
	var tadas []domain.Tada
	var query *gorm.DB = r.db.Model(&domain.Tada{}).
		Preload("Creator").Preload("Assignee").
		Where("created_by = ?", userID)

	// Set default limit
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	// Apply cursor pagination
	if pagination.Cursor != "" {
		cursor, err := dto.DecodeCursor(pagination.Cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor: %w", err)
		}

		query = query.Where("(created_at < ? OR (created_at = ? AND id < ?))",
			cursor.CreatedAt, cursor.CreatedAt, cursor.ID)
	}

	// Fetch one extra record to check if there are more pages
	query = query.Order("created_at DESC, id DESC").Limit(pagination.Limit + 1)

	if err := query.Find(&tadas).Error; err != nil {
		return nil, "", fmt.Errorf("failed to fetch tadas: %w", err)
	}

	// Generate next cursor
	var nextCursor string
	if len(tadas) > pagination.Limit {
		tadas = tadas[:pagination.Limit] // Remove extra record
		lastTada := tadas[len(tadas)-1]
		nextCursor = dto.EncodeCursor(dto.Cursor{
			ID:        lastTada.ID.String(),
			CreatedAt: lastTada.CreatedAt,
		})
	}

	return tadas, nextCursor, nil
}

func (r *tadaRepository) GetByAssigneeID(assigneeID uuid.UUID, pagination dto.PaginationQuery) ([]domain.Tada, string, error) {
	var tadas []domain.Tada
	var query *gorm.DB = r.db.Model(&domain.Tada{}).
		Preload("Creator").Preload("Assignee").
		Where("assigned_to = ?", assigneeID)

	// Set default limit
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	// Apply cursor pagination
	if pagination.Cursor != "" {
		cursor, err := dto.DecodeCursor(pagination.Cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor: %w", err)
		}

		query = query.Where("(created_at < ? OR (created_at = ? AND id < ?))",
			cursor.CreatedAt, cursor.CreatedAt, cursor.ID)
	}

	// Fetch one extra record to check if there are more pages
	query = query.Order("created_at DESC, id DESC").Limit(pagination.Limit + 1)

	if err := query.Find(&tadas).Error; err != nil {
		return nil, "", fmt.Errorf("failed to fetch tadas: %w", err)
	}

	// Generate next cursor
	var nextCursor string
	if len(tadas) > pagination.Limit {
		tadas = tadas[:pagination.Limit] // Remove extra record
		lastTada := tadas[len(tadas)-1]
		nextCursor = dto.EncodeCursor(dto.Cursor{
			ID:        lastTada.ID.String(),
			CreatedAt: lastTada.CreatedAt,
		})
	}

	return tadas, nextCursor, nil
}
