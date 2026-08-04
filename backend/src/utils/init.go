package utils

import (
	"gorm.io/gorm"
)

func Init(db *gorm.DB) {
	initSnowFlake(db)
}
