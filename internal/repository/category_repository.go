package repository

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindByID(id uint) (*entity.Category, error)
	FindByName(name string) (*entity.Category, error)
	FindByNames(names []string) ([]entity.Category, error)
	FindBySlug(slug string) (*entity.Category, error)
	Create(category *entity.Category) error
	FindAll() ([]entity.Category, error)
	Update(category *entity.Category) error
	Delete(id uint) error
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

func (r *categoryRepositoryImpl) FindByID(id uint) (*entity.Category, error) {
	var category entity.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *categoryRepositoryImpl) FindByName(name string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	return &category, err
}

func (r *categoryRepositoryImpl) FindByNames(names []string) ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Where("name IN (?)", names).Find(&categories).Error
	return categories, err
}
func (r *categoryRepositoryImpl) FindBySlug(slug string) (*entity.Category, error) {
	var category entity.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	return &category, err
}
func (r *categoryRepositoryImpl) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepositoryImpl) FindAll() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepositoryImpl) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

// Delete menghapus kategori dari database (soft delete)
func (r *categoryRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Category{}, id).Error
}
