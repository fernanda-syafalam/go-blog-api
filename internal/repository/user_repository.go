package repository

import (
	"github.com/fernanda-syafalam/backend-monitoring-notification/internal/entity"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *zerolog.Logger
}

func NewUserRepository(log *zerolog.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) FindByToken(db *gorm.DB, entity *entity.User, token string) error {
	return db.Where("token = ?", token).Find(&entity).Error
}