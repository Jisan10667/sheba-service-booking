package repository

import (
	
	"service-booking/internal/model"

	"gorm.io/gorm"
)

type BookingRepository interface {
	FindAll(page, limit int, filters map[string]interface{}) ([]model.Booking, int64, error)
	FindByID(id uint) (*model.Booking, error)
	FindByReferenceCode(referenceCode string) (*model.Booking, error)
	Create(booking *model.Booking) error
	Update(booking *model.Booking) error
	UpdateStatus(id uint, status model.BookingStatus) error
	CreateWithStatusHistory(booking *model.Booking) error
	UpdateStatusWithHistory(id uint, status model.BookingStatus, statusHistory *model.BookingStatusHistory) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

func (r *bookingRepository) FindAll(page, limit int, filters map[string]interface{}) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var count int64

	offset := (page - 1) * limit
	query := r.db.Model(&model.Booking{})

	// Apply filters
	if filters != nil {
		for key, value := range filters {
			switch key {
			case "status":
				query = query.Where("status = ?", value)
			case "user_id":
				query = query.Where("user_id = ?", value)
			case "service_id":
				query = query.Where("service_id = ?", value)
			case "start_date":
				query = query.Where("scheduled_at >= ?", value)
			case "end_date":
				query = query.Where("scheduled_at <= ?", value)
			}
		}
	}

	// Count total records
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated results with preloading
	err = query.
		Preload("Service").
		Preload("User").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	return bookings, count, err
}

func (r *bookingRepository) FindByID(id uint) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.
		Preload("Service").
		Preload("User").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&booking, id).Error
	return &booking, err
}

func (r *bookingRepository) FindByReferenceCode(referenceCode string) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.
		Where("booking_reference_code = ?", referenceCode).
		Preload("Service").
		Preload("User").
		Preload("StatusHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&booking).Error
	return &booking, err
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) CreateWithStatusHistory(booking *model.Booking) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create booking
		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		// Create initial status history
		statusHistory := &model.BookingStatusHistory{
			BookingID: booking.ID,
			Status:    booking.Status,
			IsActive:  true,
			Notes:     "Booking created",
		}
		return tx.Create(statusHistory).Error
	})
}

func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) UpdateStatus(id uint, status model.BookingStatus) error {
	return r.db.Model(&model.Booking{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *bookingRepository) UpdateStatusWithHistory(
	id uint, 
	status model.BookingStatus, 
	statusHistory *model.BookingStatusHistory,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update booking status
		if err := tx.Model(&model.Booking{}).
			Where("id = ?", id).
			Update("status", status).Error; err != nil {
			return err
		}

		// Deactivate previous active status history
		if err := tx.Model(&model.BookingStatusHistory{}).
			Where("booking_id = ?", id).
			Update("is_active", false).Error; err != nil {
			return err
		}

		// Create new status history
		return tx.Create(statusHistory).Error
	})
}