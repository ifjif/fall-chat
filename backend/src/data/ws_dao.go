package data

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"
	"zzc/fall-chat/src/database"
	"zzc/fall-chat/src/host"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type WSDao struct {
	BaseDao
}

func NewWSDao(db *gorm.DB) *WSDao {
	bs := BaseDao{db: db, rdb: database.Redis()}
	return &WSDao{
		BaseDao: bs,
	}
}

func (dao *WSDao) AddGrpcAddr2Redis(id string) {
	key := rediskey.GrpcAddr{ID: id}
	grpcAddr := host.GRPCAddr
	// 5min
	err := dao.rdb.Set(context.Background(), key.String(), grpcAddr, key.TTL()).Err()
	if err != nil {
		utils.LogError("设置%s 失败: error=%v", key, err)
	}
}

func (dao *WSDao) DelGrpcAddr2Redis(id string) {
	key := rediskey.GrpcAddr{ID: id}
	dao.rdb.Del(context.Background(), key.String())
}

func (dao *WSDao) GetGrpcAddrFromRedis(id string) (addr string, err error) {
	key := rediskey.GrpcAddr{ID: id}

	cmd := dao.rdb.Get(context.Background(), key.String())

	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		return "", cmd.Err()
	}
	addr, _ = cmd.Result()

	return addr, nil
}

func (dao *WSDao) findConversationBySingleKeyWithDB(db *gorm.DB, singleKey string) (*model.IMConversation, error) {
	// 从缓存中找，没有从数据库
	key := rediskey.Conv{SingleKey: singleKey}

	res, err := dao.rdb.Get(context.Background(), key.String()).Result()
	if err != nil && err != redis.Nil {
		// TODO:记录错误
		utils.LogError("找会话基本信息 从redis 失败：%v", err)
	}

	var conv model.IMConversation
	if err != redis.Nil {
		err := json.Unmarshal([]byte(res), &conv)
		if err == nil {
			utils.LogDebug("从缓存中找到--------")
			return &conv, nil
		}
		// TODO: 记录json错误
		utils.LogError("json反序列化错误：%v", err)
	}

	// 从数据库
	err = db.
		Where("single_key = ?", singleKey).
		Select("id, public_id,conv_type").
		Take(&conv).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	cacheInfo := map[string]any{
		"id":        conv.ID,
		"public_id": conv.PublicID,
		"conv_type": conv.ConvType,
	}
	value, err := json.Marshal(cacheInfo)
	if err != nil {
		utils.LogError("json 序列化失败：%v", err)
		return &conv, nil
	}

	err = dao.rdb.Set(context.Background(), key.String(), value, key.TTL()).Err()
	if err != nil {
		// TODO: 记录日志
		utils.LogError("设置缓存失败：%v", err)
	}

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
	key := rediskey.ImUserConvGroupID{ID: id}

	cmd := dao.rdb.SMembers(context.Background(), key.String())
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
	args = append(args, key.TTL()/time.Second)
	args = append(args, utils.ToAnySlice(groupIds)...)

	luaScript := `
		redis.call('SADD', KEYS[1], unpack(ARGV, 2, #ARGV))
		return redis.call('EXPIRE', KEYS[1], tonumber(ARGV[1]))
	`
	cmdR := dao.rdb.Eval(context.Background(), luaScript, []string{key.String()}, args...)

	if cmdR.Err() != nil {
		// TODO: redis服务错误
		utils.LogError("设置缓存失败：%v", cmdR.Err())
	}

	return groupIds, nil
}

// 找群 中所有用户ids
func (dao *WSDao) FindGroupUserIDs(groupID int64) ([]int64, error) {
	key := rediskey.GroupUserIDs{GroupID: groupID}

	res, err := dao.rdb.SMembers(context.Background(), key.String()).Result()
	if err != nil {
		utils.LogError("从redis中找 group=%s 的用户id失败：%v", err)
	}

	var result []int64
	if len(res) > 0 {
		for _, v := range res {
			iv, _ := strconv.ParseInt(v, 0, 64)
			result = append(result, iv)
		}
		return result, nil
	}

	err = dao.db.Model(model.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("user_id", &result).Error
	if err != nil {
		return nil, err
	}

	luascript := `
		redis.call('sadd', KEYS[1], unpack(ARGV, 2, #ARGV))
		return redis.call('expire', KEYS[1], tonumber(ARGV[1]))
	`

	args := []any{
		key.TTL() / time.Second,
	}
	args = append(args, utils.ToAnySlice(result)...)
	err = dao.rdb.Eval(context.Background(), luascript, []string{key.String()}, args...).Err()
	if err != nil {
		utils.LogError("设置group=%s的所有用户id到redis失败失败: error=%v", err)
	}

	return result, err
}
