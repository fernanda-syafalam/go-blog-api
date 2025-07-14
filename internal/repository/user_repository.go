package repository

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(db *gorm.DB, entity *entity.User) error
	CountByEmail(email string) (int64, error)
	FindByID(id uint) (*entity.User, error)
	FindByToken(entity *entity.User, token string) error
	FindByEmail(email string) (*entity.User, error)
	FindAll() ([]entity.User, error)
	Update(db *gorm.DB, entity *entity.User) error
	Delete(id uint) error
}

type userRepositoryImpl struct {
	db  *gorm.DB
	Log *zerolog.Logger
}

func (r *userRepositoryImpl) Create(db *gorm.DB, entity *entity.User) error {
	return db.Create(&entity).Error
}

func (r *userRepositoryImpl) CountByEmail(email string) (int64, error) {
	var total int64
	err := r.db.Model(&entity.User{}).Where("email = ?", email).Count(&total).Error
	return total, err
}

func (r *userRepositoryImpl) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepositoryImpl) FindByToken(entity *entity.User, token string) error {
	return r.db.Where("token = ?", token).Find(&entity).Error
}

func (r *userRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Model(&entity.User{}).Where("email = ?", email).First(&user).Error

	return &user, err
}

func (r *userRepositoryImpl) FindAll() ([]entity.User, error) { // <-- Implementasi metode baru
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) Update(db *gorm.DB, entity *entity.User) error {
	return db.Save(entity).Error
}

func (r *userRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}
func NewUserRepository(log *zerolog.Logger, db *gorm.DB) UserRepository {
	return &userRepositoryImpl{
		db:  db,
		Log: log,
	}
}
