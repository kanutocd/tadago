package domain
package domain

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
    Name      string         `gorm:"size:255;not null" json:"name"`
    Email     string         `gorm:"size:255;not null;uniqueIndex" json:"email"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    // Relationships
    CreatedTadas []Tada `gorm:"foreignKey:CreatedBy" json:"created_tadas,omitempty"`
    AssignedTadas []Tada `gorm:"foreignKey:AssignedTo" json:"assigned_tadas,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

func (User) TableName() string {
    return "users"
}
