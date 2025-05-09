package handler

import (
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"service-booking/internal/model"
	"service-booking/internal/service"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService}
}


func (h *BookingHandler) GetBookings(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Prepare filters
	filters := make(map[string]interface{})
	
	// Filter by status
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// Filter by user ID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			filters["user_id"] = uint(userID)
		}
	}

	// Filter by service ID
	if serviceIDStr := c.Query("service_id"); serviceIDStr != "" {
		serviceID, err := strconv.ParseUint(serviceIDStr, 10, 32)
		if err == nil {
			filters["service_id"] = uint(serviceID)
		}
	}

	// Filter by date range
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err == nil {
			filters["start_date"] = startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err == nil {
			filters["end_date"] = endDate
		}
	}

	// Fetch bookings with filters
	bookings, count, err := h.bookingService.GetBookings(page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookings,
		"meta": gin.H{
			"total":       count,
			"page":        page,
			"limit":       limit,
			"total_pages": (count + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}
	
	booking, err := h.bookingService.GetBookingByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	
	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking model.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Optional: Get user ID from context for authenticated bookings
	userID, exists := c.Get("user_id")
	if exists {
		// Convert user ID to uint
		currentUserID, ok := userID.(uint)
		if ok {
			booking.UserID = currentUserID
		}
	}

	// Create booking
	if err := h.bookingService.CreateBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) UpdateBookingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}
	
	// Struct to handle status update request
	var statusData struct {
		Status string `json:"status" binding:"required"`
		Notes  string `json:"notes,omitempty"`
	}
	
	if err := c.ShouldBindJSON(&statusData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get user ID from context (assuming JWT middleware sets the user_id)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert user ID to uint
	currentUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Update booking status
	if err := h.bookingService.UpdateBookingStatus(
		uint(id), 
		model.BookingStatus(statusData.Status), 
		currentUserID, 
		statusData.Notes,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Booking status updated successfully",
		"status":  statusData.Status,
	})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	// Get booking ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	// Get user ID from context (assuming JWT middleware sets the user_id)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert user ID to uint
	currentUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Optional: Allow cancellation reason
	var cancelRequest struct {
		Reason string `json:"reason,omitempty"`
	}
	
	// Try to bind JSON if provided, but don't require it
	_ = c.ShouldBindJSON(&cancelRequest)

	// Prepare notes
	notes := "Booking cancelled by user"
	if cancelRequest.Reason != "" {
		notes += ". Reason: " + cancelRequest.Reason
	}

	// Attempt to cancel the booking
	err = h.bookingService.CancelBooking(uint(id), currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to cancel booking: " + err.Error(),
		})
		return
	}

	// Prepare response
	response := gin.H{
		"message": "Booking cancelled successfully",
	}
	
	// Add reason to response if provided
	if cancelRequest.Reason != "" {
		response["reason"] = cancelRequest.Reason
	}

	c.JSON(http.StatusOK, response)
}

// Other existing methods...

func (h *BookingHandler) GetBookingByReferenceCode(c *gin.Context) {
	referenceCode := c.Param("code")
	
	if referenceCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking reference code"})
		return
	}
	
	booking, err := h.bookingService.GetBookingByReferenceCode(referenceCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	
	c.JSON(http.StatusOK, booking)
}

// CreateUserBooking method for authenticated users to create their own bookings
func (h *BookingHandler) CreateUserBooking(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert user ID to uint
	currentUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse booking data
	var booking model.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID from context
	booking.UserID = currentUserID

	// Create booking
	if err := h.bookingService.CreateBooking(&booking); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, booking)
}