package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/repository"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type PostUseCase interface {
	CreatePost(title, content string, authorID uint, categoryNames []string) (*entity.Post, error)
	GetPostByID(id uint) (*entity.Post, error)
	GetPostBySlug(slug string) (*entity.Post, error)
	GetAllPosts(page, limit int) ([]entity.Post, error)
	UpdatePost(id uint, title, content *string, publishedAt *time.Time, categoryNames *[]string, authorID uint) (*entity.Post, error)
	DeletePost(id uint, authorID uint) error
}

type PostUseCaseImpl struct {
	PostRepository     repository.PostRepository
	CategoryRepository repository.CategoryRepository
	validator          *validator.Validate
}

func (s *PostUseCaseImpl) CreatePost(title, content string, authorID uint, categoryNames []string) (*entity.Post, error) {

	slug := utils.GenerateSlug(title)

	existingPost, err := s.PostRepository.FindBySlug(slug)

	if err == nil && existingPost != nil {
		return nil, utils.ErrValidation("Judul postingan sudah ada (slug " + slug + " sudah terpakai), coba judul lain")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var categories []entity.Category
	if len(categoryNames) > 0 {
		categories, err = s.CategoryRepository.FindByNames(categoryNames)
		if err != nil {
			return nil, errors.New("Gagal mencari kategori: " + err.Error())
		}

		if len(categories) != len(categoryNames) {
			return nil, utils.ErrValidation("Beberapa kategori yang diberikan tidak ditemukan")
		}
	}

	post := &entity.Post{
		Title:      title,
		Slug:       slug,
		Content:    content,
		AuthorID:   authorID,
		Categories: categories,
	}

	err = s.PostRepository.Create(post)
	if err != nil {
		return nil, errors.New("Gagal menyimpan postingan ke database: " + err.Error())
	}
	return post, nil
}

func (s *PostUseCaseImpl) DeletePost(id uint, authorID uint) error {
	post, err := s.PostRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("postingan")
		}
		return errors.New("gagal menemukan postingan untuk dihapus")
	}

	if post.AuthorID != authorID {
		return utils.ErrForbidden("Anda tidak memiliki izin untuk menghapus postingan ini")
	}

	err = s.PostRepository.Delete(id)
	if err != nil {
		return errors.New("gagal menghapus postingan dari database")
	}
	return nil
}

func (s *PostUseCaseImpl) GetAllPosts(page, limit int) ([]entity.Post, error) {
	offset := (page - 1) * limit

	posts, err := s.PostRepository.FindAll(offset, limit)
	if err != nil {
		return nil, errors.New("gagal mengambil daftar postingan")
	}
	return posts, nil
}

func (s *PostUseCaseImpl) GetPostBySlug(slug string) (*entity.Post, error) {
	post, err := s.PostRepository.FindBySlug(slug)
	fmt.Println("errore", err)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("postingan")
		}
		return nil, errors.New("gagal mengambil postingan berdasarkan slug")
	}
	return post, nil
}

func (s *PostUseCaseImpl) GetPostByID(id uint) (*entity.Post, error) {
	post, err := s.PostRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("postingan")
		}
		return nil, errors.New("gagal mengambil postingan berdasarkan ID")
	}
	return post, nil
}

func (s *PostUseCaseImpl) UpdatePost(id uint, title, content *string, publishedAt *time.Time, categoryNames *[]string, authorID uint) (*entity.Post, error) {
	post, err := s.PostRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("postingan")
		}
		return nil, errors.New("gagal menemukan postingan untuk diperbarui")
	}

	// Otorisasi: Anda mungkin ingin memeriksa apakah `authorID` ini adalah pemilik post atau admin
	// Ini bisa dilakukan di handler atau service tergantung kebutuhan.
	// Contoh sederhana:
	// if post.AuthorID != authorID {
	// 	return utils.ErrForbidden("Anda tidak memiliki izin untuk memperbarui postingan ini")
	// }

	if title != nil {
		newSlug := utils.GenerateSlug(*title)
		if newSlug != post.Slug {
			existingPost, err := s.PostRepository.FindBySlug(newSlug)
			if err == nil && existingPost != nil && existingPost.ID != post.ID {
				return nil, utils.ErrValidation("Judul postingan baru sudah ada (slug " + newSlug + " sudah terpakai), coba judul lain")
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}
		post.Title = *title
		post.Slug = newSlug
	}
	// Validasi dan update konten
	if content != nil {
		post.Content = *content
	}
	if publishedAt != nil {
		post.PublishedAt = publishedAt
	}

	// Update kategori
	if categoryNames != nil {
		var categories []entity.Category
		if len(*categoryNames) > 0 {
			categories, err = s.CategoryRepository.FindByNames(*categoryNames)
			if err != nil {
				return nil, errors.New("Gagal mencari kategori untuk update: " + err.Error())
			}
			if len(categories) != len(*categoryNames) {
				return nil, utils.ErrValidation("Beberapa kategori yang diberikan untuk update tidak ditemukan")
			}
		}
		// GORM akan menangani asosiasi Many-to-Many jika Categories diatur
		// Ini akan mengganti semua kategori yang terkait sebelumnya
		post.Categories = categories
	}

	err = s.PostRepository.Update(post)
	if err != nil {
		return nil, errors.New("Gagal memperbarui postingan di database: " + err.Error())
	}
	return post, nil
}

func NewPostUseCase(postRepo repository.PostRepository, categotyRepository repository.CategoryRepository, validator *validator.Validate) PostUseCase {
	return &PostUseCaseImpl{
		PostRepository:     postRepo,
		CategoryRepository: categotyRepository,
		validator:          validator,
	}
}
