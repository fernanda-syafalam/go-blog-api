package model

import "time"

type CreatePostRequest struct {
	Title   string   `json:"title" validate:"required,min=5,max=255"`
	Content string   `json:"content" validate:"required,min=10"`
	CategoryNames    []string `json:"categoryNames" validate:"required,min=1"`
}

type UpdatePostRequest struct {
	Title         *string    `json:"title" validate:"omitempty,min=5,max=255"`
	Content       *string    `json:"content" validate:"omitempty,min=10"`
	PublishedAt   *time.Time `json:"publishedAt"`
	CategoryNames *[]string  `json:"categoryNames" validate:"omitempty,dive,min=1,max=50"` // Pointer ke slice
}
