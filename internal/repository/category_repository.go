package repository

import (
	"service-booking/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(page, limit int, filters map[string]interface{}) ([]model.Category, int64, error)
	FindByID(id uint) (*model.Category, error)
	Create(category *model.Category) error
	Update(category *model.Category) error
	Delete(id uint) error
	FindByParentCategory(parentCategoryID uint) ([]model.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindAll(page, limit int, filters map[string]interface{}) ([]model.Category, int64, error) {
	var categories []model.Category
	var count int64

	offset := (page - 1) * limit
	query := r.db.Model(&model.Category{})

	// Apply filters
	if filters != nil {
		for key, value := range filters {
			switch key {
			case "is_active":
				query = query.Where("is_active = ?", value)
			case "parent_category_id":
				query = query.Where("parent_category_id = ?", value)
			case "name":
				query = query.Where("name LIKE ?", "%"+value.(string)+"%")
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
		Preload("ParentCategory").
		Offset(offset).
		Limit(limit).
		Order("display_order").
		Find(&categories).Error

	return categories, count, err
}

func (r *categoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.
		Preload("ParentCategory").
		Preload("Services").
		First(&category, id).Error
	return &category, err
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *categoryRepository) FindByParentCategory(parentCategoryID uint) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.
		Where("parent_category_id = ?", parentCategoryID).
		Find(&categories).Error
	return categories, err
}