package usecase

import (
	"errors"

	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/repository"
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CommentUseCase interface {
	CreateComment(content string, postID, authorID uint) (*entity.Comment, error)
	GetCommentsByPostID(postID uint) ([]entity.Comment, error)
	UpdateComment(commentID, authorID uint, content string) (*entity.Comment, error)
	DeleteComment(commentID, authorID uint) error
}

type commentUseCaseImpl struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository 
	validator   *validator.Validate
}

func NewCommentUseCase(commentRepo repository.CommentRepository, postRepo repository.PostRepository, validator *validator.Validate) CommentUseCase {
	return &commentUseCaseImpl{commentRepo: commentRepo, postRepo: postRepo, validator: validator}
}

func (s *commentUseCaseImpl) CreateComment(content string, postID, authorID uint) (*entity.Comment, error) {
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("Postingan") 
		}
		return nil, errors.New("Gagal memverifikasi postingan: " + err.Error())
	}

	comment := &entity.Comment{
		Content:  content,
		PostID:   postID,
		AuthorID: authorID,
	}

	err = s.commentRepo.Create(comment)
	if err != nil {
		return nil, errors.New("Gagal menyimpan komentar: " + err.Error())
	}
	// Muat ulang komentar untuk mendapatkan Author dan Post yang sudah di-preload
	return s.commentRepo.FindByID(comment.ID)
}

// GetCommentsByPostID mengambil semua komentar untuk postingan tertentu
func (s *commentUseCaseImpl) GetCommentsByPostID(postID uint) ([]entity.Comment, error) {
	// Validasi apakah postID valid (apakah postingan ada)
	_, err := s.postRepo.FindByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("Postingan")
		}
		return nil, errors.New("Gagal memverifikasi postingan: " + err.Error())
	}

	comments, err := s.commentRepo.FindByPostID(postID)
	if err != nil {
		return nil, errors.New("Gagal mengambil komentar: " + err.Error())
	}
	return comments, nil
}

func (s *commentUseCaseImpl) UpdateComment(commentID, authorID uint, content string) (*entity.Comment, error) {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("Komentar")
		}
		return nil, errors.New("Gagal menemukan komentar: " + err.Error())
	}

	// Otorisasi: Hanya penulis komentar yang bisa mengupdate
	if comment.AuthorID != authorID {
		return nil, utils.ErrForbidden("Anda tidak memiliki izin untuk memperbarui komentar ini")
	}


	comment.Content = content
	err = s.commentRepo.Update(comment)
	if err != nil {
		return nil, errors.New("Gagal memperbarui komentar: " + err.Error())
	}
	return s.commentRepo.FindByID(comment.ID) // Muat ulang untuk Author/Post preload
}

// DeleteComment menghapus komentar
func (s *commentUseCaseImpl) DeleteComment(commentID, authorID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrNotFound("Komentar")
		}
		return errors.New("Gagal menemukan komentar: " + err.Error())
	}

	// Otorisasi: Hanya penulis komentar atau admin postingan terkait yang bisa menghapus
	// Untuk saat ini, hanya penulis komentar. Logika admin bisa ditambahkan di sini.
	if comment.AuthorID != authorID {
		return utils.ErrForbidden("Anda tidak memiliki izin untuk menghapus komentar ini")
	}

	err = s.commentRepo.Delete(commentID)
	if err != nil {
		return errors.New("Gagal menghapus komentar: " + err.Error())
	}
	return nil
}