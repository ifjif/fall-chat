package data

import (
	"context"
	"errors"
	"zzc/fall-chat/src/database"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SocialDao struct {
	BaseDao
}

func NewSocialDao(db *gorm.DB) *SocialDao {
	bd := BaseDao{db: db, rdb: database.Redis()}
	return &SocialDao{BaseDao: bd}
}

func (dao *SocialDao) SearchFriendships(id int64, targetID int64) ([]model.Friendship, error) {
	friendships := []model.Friendship{}
	err := dao.db.
		Where("(user_id = ? and friend_id = ?) or (user_id = ? and friend_id = ?)",
			targetID, id, // 对方，我
			id, targetID, // 我，对方
		).
		Find(&friendships).Error

	return friendships, err
}

func (dao *SocialDao) SearchFriendRequestPendingAndAgreeCount(id, targetID int64) (int64, error) {
	// 0:待处理 1:已同意
	var existsCount int64
	err := dao.db.
		Model(&model.FriendshipRequest{}).
		Where(`((from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?)) 
		and status in (0,1)`,
			id, targetID,
			targetID, id,
		).Count(&existsCount).Error

	return existsCount, err
}

func (dao *SocialDao) FindGroupRequestPendingAndAgreeCount(groupID, userID int64) (int64, error) {
	var existsCount int64

	err := dao.db.Model(&model.GroupRequest{}).
		Where("to_group_id = ? and from_user_id = ? and status in (0,1)", groupID, userID).
		Count(&existsCount).Error

	return existsCount, err
}

func (dao *SocialDao) FindGroupWithStatus(groupID int64) (*model.GroupInfo, error) {
	row := &model.GroupInfo{}
	row.ID = uint64(groupID)
	err := dao.db.Select("status,join_mode").Take(row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return row, nil
}

func (dao *SocialDao) FindGroupWithMemberCount(groupID int64) (*model.GroupInfoRow, error) {
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

func (dao *SocialDao) ExistsGroupMember(groupID, userID int64) (int64, error) {
	var existsCount int64
	err := dao.db.Model(&model.GroupMember{}).
		Where("group_id = ? and user_id = ?", groupID, userID).
		Count(&existsCount).Error

	return existsCount, err
}

func (dao *SocialDao) CreateOneGroupRequest(m *model.GroupRequest) error {
	return dao.db.Create(m).Error
}

func (dao *SocialDao) CreateFriendRequest(fr *model.FriendshipRequest) error {
	return dao.db.Omit("pair_key").Create(fr).Error
}

func (dao *SocialDao) SearchFriendRequestRow(id int64, filterStatus []int, limit int) ([]model.SearchFriendRequestRow, error) {
	mo := []model.SearchFriendRequestRow{}
	qc := dao.db.Table("friendship_requests as f").
		Joins("inner join users u on (u.id = f.from_user_id or u.id = f.to_user_id) and u.id != ?", id).
		Where("f.to_user_id = ? or f.from_user_id = ?", id, id).
		Order("f.created_at desc").
		Limit(limit)

	if len(filterStatus) > 0 {
		qc.Where("f.status in ?", filterStatus)
	}

	err := qc.Select("f.*, u.nickname, u.avatar").Find(&mo).Error

	return mo, err
}

func (dao *SocialDao) SearchFriendshipByToUerID(id int64, toUserID int64) (model.FriendshipRequest, error) {
	var row model.FriendshipRequest
	err := dao.db.Where("id = ? and to_user_id = ?", id, toUserID).Take(&row).Error
	return row, err
}

func (dao *SocialDao) UpdateFriendshipRequestStatusWithTx(tx *gorm.DB, targetID int64, updates map[string]any) error {
	return tx.Model(&model.FriendshipRequest{}).
		Where("id = ?", targetID).
		Updates(updates).Error
}

func (dao *SocialDao) CreateFriendshipsWithTx(tx *gorm.DB, list []model.Friendship) error {
	return tx.Create(list).Error
}

func (dao *SocialDao) SearchFriendshipList(id int64) ([]model.FriendshipRow, error) {
	var rows []model.FriendshipRow
	err := dao.db.
		Table("friendships as f").
		Joins("inner join users u on u.id = f.friend_id").
		Where("f.user_id = ?", id).
		Select("u.public_id, f.Remark, f.Source, u.nickname, u.avatar").
		Limit(10).
		Find(&rows).Error

	return rows, err
}

// NOTE: 找我管理的群 和 我请求的群:
func (dao *SocialDao) SearchGroupRequestRow(id int64, filterStatus []int) ([]model.GroupRequestRow, error) {
	// INFO: 找我管理的群
	subquery := dao.db.Model(&model.GroupMember{}).
		Where("user_id = ? and role in ?", id, []int{1, 2}).
		Select("group_id")

	query1 := dao.db.Model(&model.GroupRequest{}).
		Where("to_group_id in (?)", subquery).
		Select("id")

	// INFO: 找我的请求
	query2 := dao.db.Model(&model.GroupRequest{}).
		Where("from_user_id = ?", id).
		Select("id")

	if len(filterStatus) > 0 {
		query1.Where("status in ?", filterStatus)
		query2.Where("status in ?", filterStatus)
	}

	var ids []int64
	err := dao.db.Raw("select id from (? union all ?) as combined order by id desc limit ?",
		query1, query2, 20).
		Scan(&ids).Error
	if err != nil {
		return nil, err
	}

	var rows []model.GroupRequestRow

	if len(ids) <= 0 {
		return rows, err
	}

	err = dao.db.Model(&model.GroupRequest{}).
		Preload("GroupInfo", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, public_id, group_name, avatar")
		}).
		Preload("ApplyUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, public_id, nickname, avatar")
		}).
		Preload("InviterUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, public_id, nickname, avatar")
		}).
		Preload("HandleByUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, public_id, nickname, avatar")
		}).
		Where("id in ?", ids).
		Order("created_at desc").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (dao *SocialDao) FindGroupRequestPendingInfo(id int64) (*model.GroupRequest, error) {
	m := &model.GroupRequest{}
	m.ID = uint64(id)
	err := dao.db.
		Select("to_group_id, from_user_id, status, apply_type").
		Take(m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return m, err
}

func (dao *SocialDao) FindGroupMemberWithRole(groupID int64, userID int64) (*model.GroupMember, error) {
	m := &model.GroupMember{}
	err := dao.db.Where("group_id = ? and user_id = ?", groupID, userID).
		Select("role").
		Take(m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return m, err
}

func (dao *SocialDao) UpdateGroupRequestStatusWithTx(tx *gorm.DB, updates *model.GroupRequest) error {
	return tx.Model(updates).Updates(updates).Error
}

func (dao *SocialDao) CreateGroupMemberWithTx(tx *gorm.DB, member *model.GroupMember) error {
	return tx.Create(member).Error
}

func (dao *SocialDao) CreateGroupConvMemberWithTx(tx *gorm.DB, convMember *model.IMConversationMember) error {
	return tx.Create(convMember).Error
}

func (dao *SocialDao) FindGroupConvID(groupID string) (int64, error) {
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

func (dao *SocialDao) SetGroupCache(groupID string, userID int64) {
	key := rediskey.ImUserConvGroupID{ID: userID}

	luaScript := `
		if redis.call('exists',keys[1]) == 1 then
			return redis.call('sadd',keys[1],argv[1])
		else
			return -1
		end
	`

	cmd := dao.rdb.Eval(context.Background(), luaScript, []string{key.String()}, groupID)
	err := cmd.Err()
	if err != nil && err != redis.Nil {
		// WARN: redis 服务错误
	}

	/*
		result, _ := cmd.Result()
		-1表示键不存在
		0表示设置失败
	*/
}
