package entity

type Category struct {
	BaseEntity
	Name string `gorm:"colomn:name;not null" json:"name"`
	Slug string `json:"slug" gorm:"unique;not null"` 
}

func (*Category) TableName() string {
	return "categories"
}