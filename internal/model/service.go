package model

import (
	"time"
)

type Service struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	CategoryID  uint           `gorm:"not null" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	IsActive    bool           `json:"is_active"`
	UpdatedBy   uint           `gorm:"foreignKey:UpdatedBy;references:ID" json:"updated_by,omitempty"`
}
