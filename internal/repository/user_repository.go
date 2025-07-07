package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kanutocd/tada/internal/domain"
	"github.com/kanutocd/tada/internal/dto"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepository) GetAll(pagination dto.PaginationQuery) ([]domain.User, string, error) {
	var users []domain.User
	var query *gorm.DB = r.db.Model(&domain.User{})

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

	if err := query.Find(&users).Error; err != nil {
		return nil, "", fmt.Errorf("failed to fetch users: %w", err)
	}

	// Generate next cursor
	var nextCursor string
	if len(users) > pagination.Limit {
		users = users[:pagination.Limit] // Remove extra record
		lastUser := users[len(users)-1]
		nextCursor = dto.EncodeCursor(dto.Cursor{
			ID:        lastUser.ID.String(),
			CreatedAt: lastUser.CreatedAt,
		})
	}

	return users, nextCursor, nil
}
