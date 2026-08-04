package model

import "gorm.io/gorm"

func Init(db *gorm.DB) {
	db.AutoMigrate(&MachineIdPool{})
	db.AutoMigrate(&User{})
}
