package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"service-booking/internal/model"
	"service-booking/internal/service"
)

type ServiceHandler struct {
	serviceService service.ServiceService
}

func NewServiceHandler(serviceService service.ServiceService) *ServiceHandler {
	return &ServiceHandler{serviceService}
}

func (h *ServiceHandler) GetServices(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	categoryIDStr := c.Query("category_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	var services []model.Service
	var count int64

	if categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		
		services, count, err = h.serviceService.GetServicesByCategory(uint(categoryID), page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
			return
		}
	} else {
		services, count, err = h.serviceService.GetServices(page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": services,
		"meta": gin.H{
			"total":       count,
			"page":        page,
			"limit":       limit,
			"total_pages": (count + int64(limit) - 1) / int64(limit),
		},
	})
}

func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}
	
	service, err := h.serviceService.GetServiceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	
	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	var service model.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.serviceService.CreateService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}
	
	c.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}
	
	// Fetch existing service
	existingService, err := h.serviceService.GetServiceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}
	
	// Bind update data
	var updateData model.Service
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update fields
	updateData.ID = existingService.ID
	
	if err := h.serviceService.UpdateService(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}
	
	c.JSON(http.StatusOK, updateData)
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}
	
	if err := h.serviceService.DeleteService(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}