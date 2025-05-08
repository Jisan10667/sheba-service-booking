package service

import (
	"service-booking/internal/model"
	"service-booking/internal/repository"
)

type ServiceService interface {
	GetServices(page, limit int) ([]model.Service, int64, error)
	GetServiceByID(id uint) (*model.Service, error)
	CreateService(service *model.Service) error
	UpdateService(service *model.Service) error
	DeleteService(id uint) error
	GetServicesByCategory(categoryID uint, page, limit int) ([]model.Service, int64, error)
}

type serviceService struct {
	serviceRepo repository.ServiceRepository
}

func NewServiceService(serviceRepo repository.ServiceRepository) ServiceService {
	return &serviceService{serviceRepo}
}

func (s *serviceService) GetServices(page, limit int) ([]model.Service, int64, error) {
	return s.serviceRepo.FindAll(page, limit)
}

func (s *serviceService) GetServiceByID(id uint) (*model.Service, error) {
	return s.serviceRepo.FindByID(id)
}

func (s *serviceService) CreateService(service *model.Service) error {
	return s.serviceRepo.Create(service)
}

func (s *serviceService) UpdateService(service *model.Service) error {
	return s.serviceRepo.Update(service)
}

func (s *serviceService) DeleteService(id uint) error {
	return s.serviceRepo.Delete(id)
}

func (s *serviceService) GetServicesByCategory(categoryID uint, page, limit int) ([]model.Service, int64, error) {
	return s.serviceRepo.FindByCategory(categoryID, page, limit)
}
