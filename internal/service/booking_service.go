package service

import (
	"errors"
	"fmt"
	"time"

	"service-booking/internal/model"
	"service-booking/internal/repository"
)

type BookingService interface {
	GetBookings(page, limit int, filters map[string]interface{}) ([]model.Booking, int64, error)
	GetBookingByID(id uint) (*model.Booking, error)
	GetBookingByReferenceCode(referenceCode string) (*model.Booking, error)
	CreateBooking(booking *model.Booking) error
	UpdateBooking(booking *model.Booking) error
	UpdateBookingStatus(id uint, status model.BookingStatus, userID uint, notes string) error
	CancelBooking(id uint, userID uint) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	serviceRepo repository.ServiceRepository
	userRepo    repository.UserRepository
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	serviceRepo repository.ServiceRepository,
	userRepo repository.UserRepository,
) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		serviceRepo: serviceRepo,
		userRepo:    userRepo,
	}
}

func (s *bookingService) GetBookings(page, limit int, filters map[string]interface{}) ([]model.Booking, int64, error) {
	return s.bookingRepo.FindAll(page, limit, filters)
}

func (s *bookingService) GetBookingByID(id uint) (*model.Booking, error) {
	return s.bookingRepo.FindByID(id)
}

func (s *bookingService) GetBookingByReferenceCode(referenceCode string) (*model.Booking, error) {
	return s.bookingRepo.FindByReferenceCode(referenceCode)
}

func (s *bookingService) CreateBooking(booking *model.Booking) error {
	// Validate service
	service, err := s.serviceRepo.FindByID(booking.ServiceID)
	if err != nil {
		return errors.New("service not found")
	}

	// Validate user
	_, err = s.userRepo.FindByID(booking.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Set default status
	if booking.Status == "" {
		booking.Status = model.BookingStatusPending
	}

	// Calculate total price
	booking.TotalPrice = service.Price * float64(booking.Duration)

	// Generate booking reference code (if not provided)
	if booking.BookingReferenceCode == "" {
		booking.BookingReferenceCode = generateBookingReferenceCode()
	}

	// Create booking and initial status history
	return s.bookingRepo.CreateWithStatusHistory(booking)
}

func (s *bookingService) UpdateBooking(booking *model.Booking) error {
	// Validate service if service ID is changed
	if booking.ServiceID > 0 {
		_, err := s.serviceRepo.FindByID(booking.ServiceID)
		if err != nil {
			return errors.New("service not found")
		}
	}

	return s.bookingRepo.Update(booking)
}

func (s *bookingService) UpdateBookingStatus(
	id uint, 
	status model.BookingStatus, 
	userID uint, 
	notes string,
) error {
	// Validate status
	validStatuses := []model.BookingStatus{
		model.BookingStatusPending,
		model.BookingStatusConfirmed,
		model.BookingStatusInProgress,
		model.BookingStatusCompleted,
		model.BookingStatusCancelled,
	}

	var isValidStatus bool
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		return errors.New("invalid booking status")
	}

	// Create status history entry
	statusHistory := &model.BookingStatusHistory{
		BookingID:     id,
		Status:        status,
		IsActive:      true,
		Notes:         notes,
		CreatedBy:     userID,
		EstimatedCompletionTime: calculateEstimatedCompletionTime(status),
	}

	return s.bookingRepo.UpdateStatusWithHistory(id, status, statusHistory)
}

func (s *bookingService) CancelBooking(id uint, userID uint) error {
	return s.UpdateBookingStatus(
		id, 
		model.BookingStatusCancelled, 
		userID, 
		"Booking cancelled by user",
	)
}

// Helper functions
func generateBookingReferenceCode() string {
	// Implement a unique booking reference code generation logic
	return fmt.Sprintf("SB-%d", time.Now().UnixNano())
}

func calculateEstimatedCompletionTime(status model.BookingStatus) *time.Time {
	now := time.Now()
	var completionTime time.Time

	switch status {
	case model.BookingStatusConfirmed:
		completionTime = now.Add(24 * time.Hour)
	case model.BookingStatusInProgress:
		completionTime = now.Add(2 * time.Hour)
	case model.BookingStatusCompleted:
		completionTime = now
	default:
		return nil
	}

	return &completionTime
}