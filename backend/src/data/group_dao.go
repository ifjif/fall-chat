package data

import (
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type GroupDao struct {
	db *gorm.DB
}

func NewGroupDao(db *gorm.DB) *GroupDao {
	return &GroupDao{
		db: db,
	}
}

func (dao *GroupDao) CreateOneGroupWithTx(tx *gorm.DB, m *model.GroupInfo) error {
	return tx.Create(m).Error
}

func (dao *GroupDao) CreateOneGroupMemberWithTx(tx *gorm.DB, m *model.GroupMember) error {
	return tx.Create(m).Error
}

func (dao *GroupDao) CreateOneGropConvWithTx(tx *gorm.DB, m *model.IMConversation) error {
	return tx.Omit("last_msg_time").Create(m).Error
}

func (dao *GroupDao) CreateOneGroupConvMemberWithTx(tx *gorm.DB, m *model.IMConversationMember) error {
	return tx.Create(m).Error
}

// 分页
// 先通过member找到 total
// 然后 分页 member得到 group_id
// 然后 找group detail
func (dao *GroupDao) FindJoinedGroups(id int64, filter []int) ([]model.GroupInfo, error) {
	var ids []int64
	q := dao.db.Model(&model.GroupMember{}).
		Where("user_id = ?", id)

	if len(filter) > 0 {
		q.Where("role in ?", filter)
	}

	err := q.Select("group_id").Limit(20).Find(&ids).Error
	if err != nil {
		return nil, err
	}

	var list []model.GroupInfo
	if len(ids) > 0 {
		err = dao.db.
			Where("id in ?", ids).
			Select("public_id, group_name, description, avatar").
			Find(&list).Error
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
