package entity

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleAuthor UserRole = "Author"
	UserRoleReader UserRole = "reader"
)

type User struct {
	BaseEntity
	Email        string   `gorm:"colomn:email:unique;not null" json:"email"`
	PasswordHash string   `gorm:"colomn:password_hash;not null" json:"password_hash"`
	Username     string   `gorm:"colomn:username;not null" json:"username"`
	Role         UserRole `gorm:"colomn:role;default:reader" json:"role"`
	Posts        []Post   `json:"posts" gorm:"foreignKey:AuthorID"`
}

func (*User) TableName() string {
	return "users"
}
