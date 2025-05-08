package model

import (
	"time"

	"gorm.io/gorm"
)
const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusCompleted BookingStatus = "completed"
)
type BookingStatusHistory struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	BookingID uint           `gorm:"not null" json:"booking_id"`
	Status    BookingStatus  `gorm:"size:20;not null" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
