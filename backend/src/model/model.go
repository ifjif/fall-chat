package model

import (
	"time"
	"zzc/fall-chat/src/utils"

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

func InjectID(m *Model) error {
	id, err := utils.SF.NextID()
	if err != nil {
		return err
	}
	m.ID = uint64(id)
	m.PublicID = utils.SF.Encode(id)

	return nil
}

func EncodeID(id int64) string {
	return utils.SF.Encode(id)
}

func DecodeID(id string) (int64, error) {
	return utils.SF.Decode(id)
}
