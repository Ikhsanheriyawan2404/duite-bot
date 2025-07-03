package repository

import (
	"finance-bot/model"
	"gorm.io/gorm"
)


type CategoryRepository interface {
	GetDefaultCategories() ([]model.Category, error)
	GetDefaultCategoriesByType(typ model.CategoryType) ([]model.Category, error)
}

type categoryRepository struct {
    db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
    return &categoryRepository{db}
}

func (r *categoryRepository) GetDefaultCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("user_id IS NULL").Find(&categories).Error
	return categories, err	
}

func (r *categoryRepository) GetDefaultCategoriesByType(typ model.CategoryType) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.
		Where("user_id IS NULL AND type = ?", typ).
		Find(&categories).Error
	return categories, err
}