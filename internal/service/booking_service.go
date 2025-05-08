package service

import (
	"errors"
	"service-booking/internal/model"
	"service-booking/internal/repository"
)

type BookingService interface {
	GetBookings(page, limit int) ([]model.Booking, int64, error)
	GetBookingByID(id uint) (*model.Booking, error)
	CreateBooking(booking *model.Booking) error
	UpdateBooking(booking *model.Booking) error
	UpdateBookingStatus(id uint, status model.BookingStatus) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	serviceRepo repository.ServiceRepository
}

func NewBookingService(bookingRepo repository.BookingRepository, serviceRepo repository.ServiceRepository) BookingService {
	return &bookingService{bookingRepo, serviceRepo}
}

func (s *bookingService) GetBookings(page, limit int) ([]model.Booking, int64, error) {
	return s.bookingRepo.FindAll(page, limit)
}

func (s *bookingService) GetBookingByID(id uint) (*model.Booking, error) {
	return s.bookingRepo.FindByID(id)
}

func (s *bookingService) CreateBooking(booking *model.Booking) error {
	// Validate service exists
	_, err := s.serviceRepo.FindByID(booking.ServiceID)
	if err != nil {
		return errors.New("service not found")
	}
	
	// Default status to pending if not set
	if booking.Status == "" {
		booking.Status = model.BookingStatusPending
	}
	
	return s.bookingRepo.Create(booking)
}

func (s *bookingService) UpdateBooking(booking *model.Booking) error {
	return s.bookingRepo.Update(booking)
}

func (s *bookingService) UpdateBookingStatus(id uint, status model.BookingStatus) error {
	// Validate status value
	validStatuses := []model.BookingStatus{
		model.BookingStatusPending,
		model.BookingStatusConfirmed,
		model.BookingStatusCancelled,
		model.BookingStatusCompleted,
	}
	
	validStatus := false
	for _, s := range validStatuses {
		if s == status {
			validStatus = true
			break
		}
	}
	
	if !validStatus {
		return errors.New("invalid booking status")
	}
	
	return s.bookingRepo.UpdateStatus(id, status)
}