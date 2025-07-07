package database

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/kanutocd/tada/internal/domain"
)

func Seed(db *gorm.DB) error {
	// Check if data already exists
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	// Create users
	users := []domain.User{
		{
			ID:    uuid.New(),
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		{
			ID:    uuid.New(),
			Name:  "Jane Smith",
			Email: "jane.smith@example.com",
		},
		{
			ID:    uuid.New(),
			Name:  "Bob Wilson",
			Email: "bob.wilson@example.com",
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to seed user %s: %w", user.Name, err)
		}
	}

	// Create tadas
	dueDate := time.Now().Add(7 * 24 * time.Hour) // 1 week from now
	tadas := []domain.Tada{
		{
			ID:          uuid.New(),
			Name:        "Complete project documentation",
			Description: "Write comprehensive documentation for the Tada API project",
			CreatedBy:   users[0].ID,
			AssignedTo:  &users[1].ID,
			Status:      domain.StatusInProgress,
			DueAt:       &dueDate,
		},
		{
			ID:          uuid.New(),
			Name:        "Review pull requests",
			Description: "Review and merge pending pull requests",
			CreatedBy:   users[1].ID,
			AssignedTo:  &users[0].ID,
			Status:      domain.StatusInProgress,
		},
		{
			ID:          uuid.New(),
			Name:        "Setup CI/CD pipeline",
			Description: "Configure GitHub Actions for automated testing and deployment",
			CreatedBy:   users[0].ID,
			AssignedTo:  &users[2].ID,
			Status:      domain.StatusCompleted,
			CompletedAt: &time.Time{},
		},
		{
			ID:          uuid.New(),
			Name:        "Database optimization",
			Description: "Optimize database queries and add missing indexes",
			CreatedBy:   users[2].ID,
			Status:      domain.StatusCancelled,
		},
	}

	// Set completed time for completed task
	completedTime := time.Now().Add(-24 * time.Hour)
	tadas[2].CompletedAt = &completedTime

	for _, tada := range tadas {
		if err := db.Create(&tada).Error; err != nil {
			return fmt.Errorf("failed to seed tada %s: %w", tada.Name, err)
		}
	}

	log.Println("Database seeded successfully")
	return nil
}
