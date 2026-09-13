package utils

import (
	"strconv"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

func init() {
	color.NoColor = false
	host, ok := GetEnv("STMP_HOST")
	if !ok {
		LogWarn("未设置STMP_HOST环境变量，无法使用邮件服务")
	}
	user, ok := GetEnv("STMP_USER")
	if !ok {
		LogWarn("未设置STMP_USER环境变量，无法使用邮件服务")
	}
	port, ok := GetEnv("STMP_PORT")
	if !ok {
		LogWarn("未设置STMP_PORT环境变量，无法使用邮件服务")
	} else {
		v, err := strconv.ParseInt(port, 10, 0)
		if err != nil {
			LogWarn("STMP_PORT值不是数字，解析失败：%v", err)
		}
		STMP_PORT = int(v)
	}
	pass, ok := GetEnv("STMP_PASS")
	if !ok {
		LogWarn("未设置STMP_PASSWORD环境变量，无法使用邮件服务")
	}

	STMP_HOST = host
	STMP_USER = user
	STMP_PASS = pass
}

func Init(db *gorm.DB) {
	initSnowFlake(db)
	initLogger()
}
