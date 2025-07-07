package domain

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type TadaStatus string

const (
    StatusInProgress TadaStatus = "in_progress"
    StatusCancelled  TadaStatus = "cancelled"
    StatusCompleted  TadaStatus = "completed"
)

type Tada struct {
    ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
    Name        string         `gorm:"size:255;not null" json:"name"`
    Description string         `gorm:"type:text" json:"description"`
    CreatedBy   uuid.UUID      `gorm:"type:uuid;not null;index" json:"created_by"`
    AssignedTo  *uuid.UUID     `gorm:"type:uuid;index" json:"assigned_to"`
    Status      TadaStatus     `gorm:"type:varchar(20);not null;default:'in_progress'" json:"status"`
    DueAt       *time.Time     `json:"due_at"`
    CompletedAt *time.Time     `json:"completed_at"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

    // Relationships
    Creator  User  `gorm:"foreignKey:CreatedBy;references:ID" json:"creator,omitempty"`
    Assignee *User `gorm:"foreignKey:AssignedTo;references:ID" json:"assignee,omitempty"`
}

func (t *Tada) BeforeCreate(tx *gorm.DB) error {
    if t.ID == uuid.Nil {
        t.ID = uuid.New()
    }
    if t.Status == "" {
        t.Status = StatusInProgress
    }
    return nil
}

func (t *Tada) BeforeUpdate(tx *gorm.DB) error {
    if t.Status == StatusCompleted && t.CompletedAt == nil {
        now := time.Now()
        t.CompletedAt = &now
    }
    return nil
}

func (Tada) TableName() string {
    return "tadas"
}
