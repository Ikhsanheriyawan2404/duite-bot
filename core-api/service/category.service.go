package service

import (
	"finance-bot/model"
	"finance-bot/repository"
)

type CategoryService interface {
	GetDefaultCategories() ([]model.Category, error)
	GetDefaultCategoriesByType(typ model.CategoryType) ([]model.Category, error)
}

type categoryService struct {
    categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
    return &categoryService{categoryRepo}
}

func (s *categoryService) GetDefaultCategories() ([]model.Category, error) {
    return s.categoryRepo.GetDefaultCategories()
}

func (s *categoryService) GetDefaultCategoriesByType(typ model.CategoryType) ([]model.Category, error) {
    return s.categoryRepo.GetDefaultCategoriesByType(typ)
}