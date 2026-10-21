package dao

import (
	"context"
	"errors"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type GroupDao interface {
	// 创建一个 群
	CreateOneGroup(ctx context.Context, m *model.GroupInfo) error

	// 分页
	// 先通过member找到 total
	// 然后 分页 member得到 group_id
	// 然后 找group detail
	FindJoinedGroups(ctx context.Context, id int64, filter []int) ([]model.GroupInfo, error)

	FindStatusAndJoinMode(ctx context.Context, groupID int64) (*model.GroupInfo, error)

	FindStatusAndMembCountAndMaxMemb(ctx context.Context, groupID int64) (*model.GroupInfoRow, error)

	// 群的会话id
	FindConvID(ctx context.Context, groupID string) (int64, error)

	// 模糊匹配
	//
	// user -> 10
	// group - 10
	// user_cursor = last_user_id    inf end
	// group_cursor = last_group_id  inf end
	//
	//	already_friend	显示"发消息"	friendships 中 status=1
	//	request_pending	显示"等待验证"	friend_requests 中有待处理申请
	//	can_add	显示"添加好友"	无任何关系记录
	//	already_joined	显示"进入群聊"	group_members 中有记录
	//	need_approval	显示"申请加入"	join_mode=1 且未加入
	//	cannot_join	显示"不可加入"	join_mode=2
	//	can_join 显示加入群聊 join_mode=1
	SearchGroups(ctx context.Context, id int64, keyword string, lastGroupID, limit int64) ([]*model.UserOrGroupRelationRow, error)
}

type groupDao struct {
	db executor.Executor
}

func NewGroupDao(db executor.Executor) GroupDao {
	return &groupDao{
		db: db,
	}
}

func (dao *groupDao) CreateOneGroup(ctx context.Context, m *model.GroupInfo) error {
	return dao.db.Create(m).Error
}

// 分页
// 先通过member找到 total
// 然后 分页 member得到 group_id
// 然后 找group detail
func (dao *groupDao) FindJoinedGroups(ctx context.Context, id int64, filter []int) ([]model.GroupInfo, error) {
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

func (dao *groupDao) FindStatusAndJoinMode(ctx context.Context, groupID int64) (*model.GroupInfo, error) {
	row := &model.GroupInfo{}
	row.ID = uint64(groupID)
	err := dao.db.Select("status,join_mode").Take(row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return row, nil
}

func (dao *groupDao) FindStatusAndMembCountAndMaxMemb(ctx context.Context, groupID int64) (*model.GroupInfoRow, error) {
	subquery := dao.db.Model(&model.GroupMember{}).
		Where("group_id = ?", groupID).
		Select("count(*)")

	row := &model.GroupInfoRow{}
	err := dao.db.Model(&model.GroupInfo{}).
		Where("id = ?", groupID).
		Select("status, (?) as member_count, max_members", subquery).
		Take(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return row, err
}

func (dao *groupDao) FindConvID(ctx context.Context, groupID string) (int64, error) {
	var convID int64
	err := dao.db.Model(model.IMConversation{}).
		Where("single_key = ?", groupID).
		Select("id").
		First(&convID).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return convID, nil
	}

	return convID, err
}

func (dao *groupDao) SearchGroups(ctx context.Context, id int64, keyword string, lastGroupID, limit int64) ([]*model.UserOrGroupRelationRow, error) {
	var groupRows []*model.UserOrGroupRelationRow

	err := dao.db.Raw(`
		select  g.public_id as id,
						g.group_name as name,
						g.avatar as avatar,
						g.description as description,
						g.join_mode as join_mode,
						gr.status as request_status,
						'' as nickname,
						gm.user_id as group_user_id,
						'group' as relation_type
		from group_infos g
			left join group_members gm on
			g.id = gm.group_id and gm.user_id = ?
			left join group_requests gr on
			g.id = gr.to_group_id and gr.status = 0 and gr.from_user_id = ?
		where g.search_flag = 1 
			and g.group_name like ?
			and g.id > ?
		order by g.id asc
		limit ?
		`, id, id, keyword, lastGroupID, limit).
		Scan(&groupRows).Error
	if err != nil {
		return nil, err
	}

	return groupRows, nil
}
