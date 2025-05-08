package model

import (
	"time"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Email     string         `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Role      UserRole       `gorm:"size:20;not null;default:user" json:"role"`
	Token	 string         `gorm:"size:255" json:"token"`
	RefreshToken string      `gorm:"size:255" json:"refresh_token"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}