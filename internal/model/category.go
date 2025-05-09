package model

import (
	"time"
)

type Category struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Name             string     `gorm:"size:100;not null;unique" json:"name"`
	Description      string     `gorm:"type:text" json:"description"`
	ParentCategoryID *uint      `json:"parent_category_id"`
	ParentCategory   *Category  `gorm:"foreignKey:ParentCategoryID" json:"parent_category,omitempty"`
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	IconURL          string     `json:"icon_url"`
	DisplayOrder     int        `json:"display_order"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	Services         []Service  `gorm:"foreignKey:CategoryID" json:"services,omitempty"`
}