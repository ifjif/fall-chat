package model

type User struct {
	Model
	Username string `gorm:"size:50;uniqueIndex"`
	Nickname string `gorm:"size:50"`
	Email    string `gorm:"size:100"`
	Phone    string `gorm:"size:20"`
	Password string `gorm:"size:255"`
	Status   int8   `gorm:"default:1"` // 0：禁用，1：正常
}

func NewUser(id int64) *User {
	user := &User{}
	user.ID = uint64(id)

	return user
}
