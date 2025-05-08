package model

import (
	"time"
)

type BookingStatus string



type Booking struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ServiceID   uint           `gorm:"not null" json:"service_id"`
	Service     Service        `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	UserName    string         `gorm:"size:255;not null" json:"user_name"`
	PhoneNumber string         `gorm:"size:20;not null" json:"phone_number"`
	Email       string         `gorm:"size:255" json:"email"`
	Status      BookingStatus  `gorm:"size:20;not null;default:pending" json:"status"`
	ScheduledAt *time.Time     `json:"scheduled_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relationship to status history
	StatusHistory []BookingStatusHistory `gorm:"foreignKey:BookingID" json:"status_history,omitempty"`
}
