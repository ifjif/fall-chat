package service

import (
	"zzc/fall-chat/src/utils"

	"gorm.io/gorm"
)

var GlobalService map[string]Service = map[string]Service{}

const (
	USER_SERVICE_NAME = "UserService"
)

func Init(db *gorm.DB) {
	register(USER_SERVICE_NAME, NewUserService(USER_SERVICE_NAME, db, utils.SF))
}

func register(name string, service Service) {
	GlobalService[name] = service
}
