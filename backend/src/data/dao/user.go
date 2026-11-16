package dao

import (
	"context"
	"errors"
	"zzc/fall-chat/src/data/executor"
	"zzc/fall-chat/src/model"

	"gorm.io/gorm"
)

type UserDao interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)

	FindNicknameAndAvatar(ctx context.Context, id int64) (*model.User, error)

	FindByUsername(ctx context.Context, username string) (*model.User, error)

	CountByUsername(ctx context.Context, username string) (int64, error)

	CreateOne(ctx context.Context, user model.User) error

	UpdateOne(ctx context.Context, user model.User) error

	UpdatePassword(ctx context.Context, id int64, newPass string) error

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
	SearchUsers(ctx context.Context, id int64, keyword string, lastUserID, limit int64) ([]*model.UserOrGroupRelationRow, error)
}

type userDao struct {
	db executor.Executor
}

func NewUserDao(db executor.Executor) UserDao {
	return &userDao{
		db: db,
	}
}

func (dao *userDao) FindNicknameAndAvatar(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	user.ID = uint64(id)
	err := dao.db.
		Select("id, nickname, avatar").
		Take(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return user, err
}

func (dao *userDao) FindByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	user.ID = uint64(id)

	err := dao.db.Take(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (dao *userDao) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := dao.db.Where("username = ?", username).Take(user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (dao *userDao) CountByUsername(ctx context.Context, username string) (int64, error) {
	var count int64
	err := dao.db.Model(&model.User{}).
		Where("username = ?", username).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dao *userDao) CreateOne(ctx context.Context, user model.User) error {
	return dao.db.Create(&user).Error
}

func (dao *userDao) UpdateOne(ctx context.Context, user model.User) error {
	return dao.db.Model(&user).Updates(user).Error
}

func (dao *userDao) UpdatePassword(ctx context.Context, id int64, newPass string) error {
	var user model.User
	user.ID = uint64(id)
	return dao.db.Model(&user).Update("password", newPass).Error
}

func (dao *userDao) SearchUsers(ctx context.Context, id int64, keyword string, lastUserID, limit int64) ([]*model.UserOrGroupRelationRow, error) {
	var userRows []*model.UserOrGroupRelationRow
	err := dao.db.Raw(`
		select u.public_id as id,
					 u.username as name,
					 u.nickname as nickname,
					 u.avatar as avatar,
					 '' as description,
					 fs.status as friend_status,
					 fsq.status as request_status,
					 'user' as relation_type
		from users u
		left join friendships fs
				on u.id = fs.friend_id and fs.status = 1 and fs.user_id = ?
		left join friendship_requests fsq
				on u.id = fsq.to_user_id and fsq.status = 0 and fsq.from_user_id = ?
		where u.search_flag = 1
			and u.id != ?
			and (u.username like ? or u.nickname like ?)
			and u.id > ?
		order by u.id asc
		limit ?
		`, id, id, id, keyword, keyword, lastUserID, limit).
		Scan(&userRows).Error
	if err != nil {
		return nil, err
	}

	return userRows, nil
}
