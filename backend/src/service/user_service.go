package service

import (
	"log"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"

	"gorm.io/gorm"
)

type UserService struct {
	name string
	sf   *utils.SnowFlake
	db   *gorm.DB
}

func NewUserService(name string, db *gorm.DB, sf *utils.SnowFlake) *UserService {
	return &UserService{
		name: name,
		db:   db,
		sf:   sf,
	}
}
func (us *UserService) service() {}

func (us *UserService) generateID() (int64, error) {
	id, err := us.sf.NextID()
	if err != nil {
		log.Printf("%s: 生成Id失败：%v", us.name, err)
		return 0, err
	}

	return id, nil
}

func (us *UserService) RegisterUser(username string, password string) (string, error) {
	id, err := us.generateID()
	if err != nil {
		return "", err
	}

	pID := us.sf.Encode(id)

	user := model.NewUser(id)
	user.PublicID = pID
	user.Username = username + pID
	user.Username = "user" + pID
	user.Password = utils.EncryptPassword(password)

	us.db.Create(user)

	return pID, err
}
