package model

/*
*  id bigint primary key
*  public_id varchar(32) unique
*  created_at datetime
*  updated_at datetime
*  delete_at datetime
*  username varchar(50) unique
*  nickname varchar(50) index
*  email varchar(100)
*  avatar varchar(255)
*  phone varchar(20)
*  password varchar(255)
*  status tinyint default 1 // 0:禁用 1:正常
*  search_flag tinyint default 1 // 1:所有人可见 0:仅通过id
*
*  index idx_public_id (public_id)
*  index idx_username (username)
*
 */
type User struct {
	Model
	Username   string `gorm:"size:50;uniqueIndex"`
	Nickname   string `gorm:"size:50;index"`
	Email      string `gorm:"size:100"`
	Avatar     string `gorm:"size:255"`
	Phone      string `gorm:"size:20"`
	Password   string `gorm:"size:255"`
	Status     *int8  `gorm:"default:1"` // 0：禁用，1：正常
	SearchFlag *int8  `gorm:"default:1"` // 1:所有人可见 0:仅通过id
}

func NewUser(id int64) *User {
	user := &User{}
	user.ID = uint64(id)

	return user
}
