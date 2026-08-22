package model

import "gorm.io/gorm"

func Init(db *gorm.DB) {
	db.AutoMigrate(&MachineIdPool{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&PasswordRestCode{})
	db.AutoMigrate(&Friendship{})
	db.AutoMigrate(&FriendshipRequest{})
	db.AutoMigrate(&GroupInfo{})
	db.AutoMigrate(&GroupMember{})
	db.AutoMigrate(&FileConfig{})
}
