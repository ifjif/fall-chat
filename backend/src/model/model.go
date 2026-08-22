package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Model struct {
	ID        uint64 `gorm:"primaryKey"`
	PublicID  string `gorm:"size:32;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"` 参与 uniqueIndex会有问题，存在为null
	DeletedAt soft_delete.DeletedAt `gorm:"index"` // 存在为0
}
