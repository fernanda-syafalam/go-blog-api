package entity

import "time"

type Post struct {
	BaseEntity
	Title       string     `gorm:"colomn:title;not null" json:"title"`
	Slug        string     `gorm:"colomn:slug;not null" json:"slug"`
	Content     string     `gorm:"type:text;colomn:content" json:"content"`
	AuthorID    uint       `gorm:"colomn:author_id;not null" json:"authorId"`
	Author      User       `gorm:"foreignKey:AuthorID" json:"author"`
	PublishedAt *time.Time `gorm:"colomn:published_at" json:"publishedAt"`
	Categories  []Category `json:"categories" gorm:"many2many:post_categories;"`
	Comments    []Comment  `json:"comments"` 

}

func (*Post) TableName() string {
	return "posts"
}
