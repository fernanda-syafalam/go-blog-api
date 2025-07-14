package entity

type Comment struct {
	BaseEntity
	Content  string `json:"content" gorm:"type:text;not null"`
	PostID   uint   `json:"postId" gorm:"not null"`
	Post     Post   `json:"post"` // Relasi Belongs To Post
	AuthorID uint   `json:"authorId" gorm:"not null"`
	Author   User   `json:"author"` // Relasi Belongs To User
}