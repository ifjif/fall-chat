package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Model struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement:false" json:"id"`
	PublicID  string    `gorm:"size:32;uniqueIndex" json:"public_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"index"` 参与 uniqueIndex会有问题，存在为null
	DeletedAt soft_delete.DeletedAt `gorm:"index" json:"deleted_at"` // 存在为0
}
