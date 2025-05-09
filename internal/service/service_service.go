package service

import (
	"errors"
	"service-booking/internal/model"
	"service-booking/internal/repository"
)

type ServiceService interface {
	GetServices(page, limit int, filters map[string]interface{}) ([]model.Service, int64, error)
	GetServiceByID(id uint) (*model.Service, error)
	CreateService(service *model.Service) error
	UpdateService(service *model.Service) error
	DeleteService(id uint) error
	GetServicesByCategory(categoryID uint, page, limit int) ([]model.Service, int64, error)
	GetFeaturedServices(limit int) ([]model.Service, error)
}

type serviceService struct {
	serviceRepo repository.ServiceRepository
	categoryRepo repository.CategoryRepository
}

func NewServiceService(
	serviceRepo repository.ServiceRepository, 
	categoryRepo repository.CategoryRepository,
) ServiceService {
	return &serviceService{
		serviceRepo: serviceRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *serviceService) GetServices(page, limit int, filters map[string]interface{}) ([]model.Service, int64, error) {
	return s.serviceRepo.FindAll(page, limit, filters)
}

func (s *serviceService) GetServiceByID(id uint) (*model.Service, error) {
	return s.serviceRepo.FindByID(id)
}

func (s *serviceService) CreateService(service *model.Service) error {
	// Validate category
	_, err := s.categoryRepo.FindByID(service.CategoryID)
	if err != nil {
		return errors.New("invalid category")
	}

	// Set default values
	if service.IsActive == false {
		service.IsActive = true
	}

	return s.serviceRepo.Create(service)
}

func (s *serviceService) UpdateService(service *model.Service) error {
	// Validate category if changed
	if service.CategoryID > 0 {
		_, err := s.categoryRepo.FindByID(service.CategoryID)
		if err != nil {
			return errors.New("invalid category")
		}
	}

	return s.serviceRepo.Update(service)
}

func (s *serviceService) DeleteService(id uint) error {
	// Optional: Check if service has any active bookings before deletion
	return s.serviceRepo.Delete(id)
}

func (s *serviceService) GetServicesByCategory(categoryID uint, page, limit int) ([]model.Service, int64, error) {
	// Validate category
	_, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return nil, 0, errors.New("invalid category")
	}

	return s.serviceRepo.FindByCategory(categoryID, page, limit)
}

func (s *serviceService) GetFeaturedServices(limit int) ([]model.Service, error) {
	return s.serviceRepo.FindFeaturedServices(limit)
}