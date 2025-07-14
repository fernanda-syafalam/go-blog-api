package usecase

import (
	"errors"

	"strings"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/repository"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(name string) (*entity.Category, error)
	GetAllCategories() ([]entity.Category, error)
	GetCategoryByID(id uint) (*entity.Category, error)
	UpdateCategory(id uint, name string) (*entity.Category, error)
	DeleteCategory(id uint) error
}

type categoryServiceImpl struct {
	categoryRepo repository.CategoryRepository
	validator    *validator.Validate
}

func NewCategoryUseCase(categoryRepo repository.CategoryRepository, validator *validator.Validate) CategoryService {
	return &categoryServiceImpl{categoryRepo: categoryRepo, validator: validator}
}

func (s *categoryServiceImpl) CreateCategory(name string) (*entity.Category, error) {
	slug := utils.GenerateSlug(name)

	existingCategoryByName, err := s.categoryRepo.FindByName(name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("Gagal memeriksa nama kategori: " + err.Error())
	}
	if existingCategoryByName != nil && existingCategoryByName.ID != 0 {
		return nil, utils.ErrValidation("Nama kategori '" + name + "' sudah ada")
	}

	existingCategoryBySlug, err := s.categoryRepo.FindBySlug(slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("Gagal memeriksa slug kategori: " + err.Error())
	}
	if existingCategoryBySlug != nil && existingCategoryBySlug.ID != 0 {
		return nil, utils.ErrValidation("Slug kategori '" + slug + "' sudah ada, coba nama lain")
	}

	category := &entity.Category{
		Name: name,
		Slug: slug,
	}

	err = s.categoryRepo.Create(category)
	if err != nil {
		return nil, errors.New("Gagal menyimpan kategori: " + err.Error())
	}
	return category, nil
}

// GetAllCategories mengambil semua kategori
func (s *categoryServiceImpl) GetAllCategories() ([]entity.Category, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, errors.New("Gagal mengambil semua kategori: " + err.Error())
	}
	return categories, nil
}

// GetCategoryByID mengambil kategori berdasarkan ID
func (s *categoryServiceImpl) GetCategoryByID(id uint) (*entity.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("Kategori")
		}
		return nil, errors.New("Gagal mengambil kategori: " + err.Error())
	}
	return category, nil
}

func (s *categoryServiceImpl) UpdateCategory(id uint, name string) (*entity.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("Kategori")
		}
		return nil, errors.New("Gagal menemukan kategori: " + err.Error())
	}

	if strings.EqualFold(category.Name, name) {
		return category, nil
	}

	existingCategoryByName, err := s.categoryRepo.FindByName(name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("Gagal memeriksa nama kategori: " + err.Error())
	}
	if existingCategoryByName != nil && existingCategoryByName.ID != 0 && existingCategoryByName.ID != id {
		return nil, utils.ErrValidation("Nama kategori '" + name + "' sudah digunakan oleh kategori lain")
	}

	category.Name = name
	category.Slug = utils.GenerateSlug(name)

	existingCategoryBySlug, err := s.categoryRepo.FindBySlug(category.Slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("Gagal memeriksa slug kategori: " + err.Error())
	}
	if existingCategoryBySlug != nil && existingCategoryBySlug.ID != 0 && existingCategoryBySlug.ID != id {
		return nil, utils.ErrValidation("Slug kategori '" + category.Slug + "' sudah digunakan oleh kategori lain, coba nama lain")
	}

	err = s.categoryRepo.Update(category)
	if err != nil {
		return nil, errors.New("Gagal memperbarui kategori: " + err.Error())
	}
	return category, nil
}

// DeleteCategory menghapus kategori
func (s *categoryServiceImpl) DeleteCategory(id uint) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("Kategori")
		}
		return errors.New("Gagal menemukan kategori: " + err.Error())
	}

	err = s.categoryRepo.Delete(id)
	if err != nil {
		return errors.New("Gagal menghapus kategori: " + err.Error())
	}
	return nil
}
