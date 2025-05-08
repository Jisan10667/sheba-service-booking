package handler

import (
	"net/http"
	"strconv"

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

	bookings, count, err := h.bookingService.GetBookings(page, limit)
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
	
	var statusData struct {
		Status string `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&statusData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.bookingService.UpdateBookingStatus(uint(id), model.BookingStatus(statusData.Status)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated successfully"})
}