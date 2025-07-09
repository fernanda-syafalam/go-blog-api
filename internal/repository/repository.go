package repository

import (
	"context"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(&entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(&entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(&entity).Error
}

func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindById(ctx context.Context, db *gorm.DB, entity *T, id any) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "Repository.FindById")
	defer span.End()
	return db.Select("id", "name", "password").Where("id = ?", id).First(&entity).Error
}

