package entity

import "time"

type BaseEntity struct {
	ID        uint       `gorm:"primaryKey"`
	CreatedAt *time.Time `gorm:"colomn:created_at"`
	UpdatedAt *time.Time `gorm:"colomn:updated_at"`
}
