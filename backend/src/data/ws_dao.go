package data

import (
	"context"
	"errors"
	"fmt"
	"time"
	"zzc/fall-chat/src/database"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ConvCache map[string]*model.IMConversation

type WSDao struct {
	BaseDao
	convCache ConvCache
}

func NewWSDao(db *gorm.DB) *WSDao {
	bs := BaseDao{db: db, rdb: database.Redis()}
	return &WSDao{
		BaseDao:   bs,
		convCache: make(ConvCache),
	}
}

func (dao *WSDao) findConversationBySingleKeyWithDB(db *gorm.DB, singleKey string) (*model.IMConversation, error) {
	// 从缓存中找，没有从数据库
	if conv, ok := dao.convCache[singleKey]; ok {
		return conv, nil
	}

	var conv model.IMConversation
	err := db.
		Where("single_key = ?", singleKey).
		Select("id, public_id,conv_type").
		Take(&conv).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// 设置缓存
	dao.convCache[singleKey] = &conv

	return &conv, nil
}

func (dao *WSDao) FindConversationBySingleKeyWithTx(tx *gorm.DB, singleKey string) (*model.IMConversation, error) {
	return dao.findConversationBySingleKeyWithDB(tx, singleKey)
}

func (dao *WSDao) FindConversationBySingleKey(singleKey string) (*model.IMConversation, error) {
	return dao.findConversationBySingleKeyWithDB(dao.db, singleKey)
}

func (dao *WSDao) CreateConversationWithTx(tx *gorm.DB, m *model.IMConversation) (bool, error) {
	err := tx.Omit("last_msg_time").Create(m).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true, nil
	}

	return false, err
}

func (dao *WSDao) CreateConversationMemberBatchWithTx(tx *gorm.DB, mms []*model.IMConversationMember) error {
	return tx.Create(mms).Error
}

func (dao *WSDao) UpdateConversation(m *model.IMConversation) error {
	return dao.db.Updates(m).Error
}

func (dao *WSDao) CreateIMMesage(m *model.IMMessage) error {
	return dao.db.Create(m).Error
}

func (dao *WSDao) FindUserProfile(id int64) (*model.User, error) {
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

func (dao *WSDao) FindGroupConvIDsByUserID(id int64) ([]string, error) {
	// 先从redis
	key := fmt.Sprintf("im:user:conv:groupid:%d", id)

	cmd := dao.rdb.SMembers(context.Background(), key)
	err := cmd.Err()
	if err != nil && err != redis.Nil {
		// TODO: 记录 redis服务错误
	}

	result, _ := cmd.Result()

	if len(result) > 0 {
		return result, nil
	}

	var groupIds []string
	err = dao.db.Model(model.IMConversationMember{}).
		Table("im_conversation_members icm").
		Joins("inner join im_conversations ic on ic.id = icm.conv_id and ic.conv_type = ? and ic.deleted_at = 0", 1).
		Joins("inner join group_infos gi on gi.id = icm.peer_id and gi.deleted_at = 0").
		Where("icm.user_id = ?", id).
		Pluck("gi.public_id", &groupIds).Error
	if err != nil {
		return nil, err
	}

	var args []any
	args = append(args, 8*time.Hour/time.Second)
	args = append(args, utils.ToAnySlice(groupIds)...)

	luaScript := `
		redis.call('SADD', KEYS[1], unpack(ARGV, 2, #ARGV))
		return redis.call('EXPIRE', KEYS[1], tonumber(ARGV[1]))
	`
	cmdR := dao.rdb.Eval(context.Background(), luaScript, []string{key}, args...)

	if cmdR.Err() != nil {
		// TODO: redis服务错误
		utils.LogError("设置缓存失败：%v", cmdR.Err())
	}

	return groupIds, nil
}
