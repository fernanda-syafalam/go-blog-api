package repository

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *entity.Comment) error
	FindByPostID(postID uint) ([]entity.Comment, error)
	FindByID(id uint) (*entity.Comment, error)
	Update(comment *entity.Comment) error
	Delete(id uint) error
}

type commentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepositoryImpl{db: db}
}

func (r *commentRepositoryImpl) Create(comment *entity.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepositoryImpl) FindByPostID(postID uint) ([]entity.Comment, error) {
	var comments []entity.Comment
	err := r.db.Where("post_id = ?", postID).Order("created_at ASC").Preload("Author").Find(&comments).Error
	return comments, err
}

func (r *commentRepositoryImpl) FindByID(id uint) (*entity.Comment, error) {
	var comment entity.Comment
	err := r.db.Preload("Author").Preload("Post").First(&comment, id).Error
	return &comment, err
}

func (r *commentRepositoryImpl) Update(comment *entity.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Comment{}, id).Error
}