package repository

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"gorm.io/gorm"
)


type PostRepository interface {
	Create(post *entity.Post) error
	FindByID(id uint) (*entity.Post, error)
	FindBySlug(slug string) (*entity.Post, error)
	FindAll(offset, limit int) ([]entity.Post, error)
	Update(post *entity.Post) error
	Delete(id uint) error
}

type PostRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) Create(post *entity.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepositoryImpl) FindByID(id uint) (*entity.Post, error) {
	var post entity.Post
	err := r.db.Preload("Author").Preload("Categories").First(&post, id).Error
	return &post, err
}

func (r *PostRepositoryImpl) FindBySlug(slug string) (*entity.Post, error) {
	var post entity.Post
	err := r.db.Preload("Author").Preload("Categories").Where("slug = ?", slug).First(&post).Error
	return &post, err
}

func (r *PostRepositoryImpl) FindAll(offset, limit int) ([]entity.Post, error) {
	var posts []entity.Post
	// Tambahkan Preload("Categories")
	err := r.db.Offset(offset).Limit(limit).Order("published_at desc").Preload("Author").Preload("Categories").Find(&posts).Error
	return posts, err
}

func (r *PostRepositoryImpl) Update(post *entity.Post) error {
	// Untuk many-to-many, kita perlu menggunakan Asociation() atau Select() pada update
	// GORM otomatis akan menangani update relasi many2many jika post.Categories dimodifikasi
	// dan kemudian Save dipanggil. Jika ada perubahan pada asosiasi, GORM akan mengurus
	// penghapusan dan penambahan entri di tabel pivot.
	return r.db.Save(post).Error
}

// Delete menghapus postingan dari database (soft delete)
func (r *PostRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.Post{}, id).Error
}