package service

import (
	"errors"
	"service-booking/internal/model"
	"service-booking/internal/repository"
)

type CategoryService interface {
	GetCategories(page, limit int, filters map[string]interface{}) ([]model.Category, int64, error)
	GetCategoryByID(id uint) (*model.Category, error)
	CreateCategory(category *model.Category) error
	UpdateCategory(category *model.Category) error
	DeleteCategory(id uint) error
	GetSubCategories(parentCategoryID uint) ([]model.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) GetCategories(page, limit int, filters map[string]interface{}) ([]model.Category, int64, error) {
	return s.categoryRepo.FindAll(page, limit, filters)
}

func (s *categoryService) GetCategoryByID(id uint) (*model.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *categoryService) CreateCategory(category *model.Category) error {
	// Validate parent category if specified
	if category.ParentCategoryID != nil {
		_, err := s.categoryRepo.FindByID(*category.ParentCategoryID)
		if err != nil {
			return errors.New("invalid parent category")
		}
	}

	// Set default values
	if category.IsActive == false {
		category.IsActive = true
	}

	return s.categoryRepo.Create(category)
}

func (s *categoryService) UpdateCategory(category *model.Category) error {
	// Validate parent category if changed
	if category.ParentCategoryID != nil {
		_, err := s.categoryRepo.FindByID(*category.ParentCategoryID)
		if err != nil {
			return errors.New("invalid parent category")
		}
	}

	return s.categoryRepo.Update(category)
}

func (s *categoryService) DeleteCategory(id uint) error {
	// Check for existing subcategories
	subCategories, err := s.categoryRepo.FindByParentCategory(id)
	if err != nil {
		return err
	}

	if len(subCategories) > 0 {
		return errors.New("cannot delete category with subcategories")
	}

	return s.categoryRepo.Delete(id)
}

func (s *categoryService) GetSubCategories(parentCategoryID uint) ([]model.Category, error) {
	return s.categoryRepo.FindByParentCategory(parentCategoryID)
}