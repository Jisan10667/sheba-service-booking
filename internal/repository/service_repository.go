package repository

import (
	"service-booking/internal/model"

	"gorm.io/gorm"
)

type ServiceRepository interface {
	FindAll(page, limit int) ([]model.Service, int64, error)
	FindByID(id uint) (*model.Service, error)
	Create(service *model.Service) error
	Update(service *model.Service) error
	Delete(id uint) error
	FindByCategory(categoryID uint, page, limit int) ([]model.Service, int64, error)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db}
}

func (r *serviceRepository) FindAll(page, limit int) ([]model.Service, int64, error) {
	var services []model.Service
	var count int64

	offset := (page - 1) * limit

	err := r.db.Model(&model.Service{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Category").Offset(offset).Limit(limit).Find(&services).Error
	return services, count, err
}

func (r *serviceRepository) FindByID(id uint) (*model.Service, error) {
	var service model.Service
	err := r.db.Preload("Category").First(&service, id).Error
	return &service, err
}

func (r *serviceRepository) Create(service *model.Service) error {
	return r.db.Create(service).Error
}

func (r *serviceRepository) Update(service *model.Service) error {
	return r.db.Save(service).Error
}

func (r *serviceRepository) Delete(id uint) error {
	return r.db.Delete(&model.Service{}, id).Error
}

func (r *serviceRepository) FindByCategory(categoryID uint, page, limit int) ([]model.Service, int64, error) {
	var services []model.Service
	var count int64

	offset := (page - 1) * limit

	err := r.db.Model(&model.Service{}).Where("category_id = ?", categoryID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Category").Where("category_id = ?", categoryID).Offset(offset).Limit(limit).Find(&services).Error
	return services, count, err
}