package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint64 `gorm:"primaryKey"`
	PublicID  string `gorm:"size:32;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
