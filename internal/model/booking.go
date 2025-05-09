package model

import (
	"time"
)

type BookingStatus string

const (
	BookingStatusPending     BookingStatus = "pending"
	BookingStatusConfirmed   BookingStatus = "confirmed"
	BookingStatusInProgress  BookingStatus = "in_progress"
	BookingStatusCompleted   BookingStatus = "completed"
	BookingStatusCancelled   BookingStatus = "cancelled"
)

type Booking struct {
	ID                   uint                 `gorm:"primaryKey" json:"id"`
	ServiceID            uint                 `gorm:"not null" json:"service_id"`
	Service              Service              `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	UserID               uint                 `gorm:"not null" json:"user_id"`
	User                 User                 `gorm:"foreignKey:UserID" json:"user,omitempty"`
	UserName             string               `gorm:"size:255;not null" json:"user_name"`
	PhoneNumber          string               `gorm:"size:20;not null" json:"phone_number"`
	Email                string               `gorm:"size:255" json:"email"`
	Status               BookingStatus        `gorm:"size:20;not null;default:pending" json:"status"`
	ScheduledAt          *time.Time           `json:"scheduled_at,omitempty"`
	BookingReferenceCode string               `gorm:"size:50;unique" json:"booking_reference_code"`
	TotalPrice           float64              `gorm:"not null" json:"total_price"`
	Duration             int                  `gorm:"default:1" json:"duration"`
	Notes                string               `gorm:"type:text" json:"notes"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
	StatusHistory        []BookingStatusHistory `gorm:"foreignKey:BookingID" json:"status_history,omitempty"`
}