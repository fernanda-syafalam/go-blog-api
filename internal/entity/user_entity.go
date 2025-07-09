package entity

type User struct {
	ID string 	`gorm:"colomn:id;primaryKey"`
	Password string `gorm:"colomn:password"`
	Name string `gorm:"colomn:name"`
	Token string `gorm:"colomn:token"`
	CreatedAt int64 `gorm:"colomn:created_at"`
	UpdatedAt int64 `gorm:"colomn:updated_at"`
}

func (*User) TableName() string {
	return "users"
}

