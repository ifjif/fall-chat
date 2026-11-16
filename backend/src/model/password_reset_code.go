package model

import "time"

type PasswordRestCode struct {
	Model
	Code     string    `gorm:"size:10;not null"`
	ExpireAt time.Time `gorm:"not null"`
	Used     bool      `gorm:"default:false"`
}
