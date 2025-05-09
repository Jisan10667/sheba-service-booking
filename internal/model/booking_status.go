package model

import (
	"time"
)

type BookingStatusHistory struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	BookingID              uint           `gorm:"not null" json:"booking_id"`
	Booking                Booking        `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Status                 BookingStatus  `gorm:"size:20;not null" json:"status"`
	CreatedAt              time.Time      `json:"created_at"`
	IsActive               bool           `json:"is_active"`
	Notes                  string         `gorm:"type:text" json:"notes"`
	CreatedBy              uint           `json:"created_by"`
	Creator                User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	EstimatedCompletionTime *time.Time    `json:"estimated_completion_time"`
}