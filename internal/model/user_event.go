package model

import "time"

type UserEvent struct {
	ID uint `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (u *UserEvent) GetId() uint {
	return u.ID
}
