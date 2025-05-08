package repository

import (
	"service-booking/internal/model"

	"gorm.io/gorm"
)

type BookingRepository interface {
	FindAll(page, limit int) ([]model.Booking, int64, error)
	FindByID(id uint) (*model.Booking, error)
	Create(booking *model.Booking) error
	Update(booking *model.Booking) error
	UpdateStatus(id uint, status model.BookingStatus) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

func (r *bookingRepository) FindAll(page, limit int) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var count int64

	offset := (page - 1) * limit

	err := r.db.Model(&model.Booking{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Service").Offset(offset).Limit(limit).Find(&bookings).Error
	return bookings, count, err
}

func (r *bookingRepository) FindByID(id uint) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("Service").First(&booking, id).Error
	return &booking, err
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) UpdateStatus(id uint, status model.BookingStatus) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", id).Update("status", status).Error
}
