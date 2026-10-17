package data

import (
	"context"
	"zzc/fall-chat/src/database"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MessageSeqDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewMessageSeqDao(db *gorm.DB) *MessageSeqDao {
	return &MessageSeqDao{
		db:  db,
		rdb: database.Redis(),
	}
}

func (dao *MessageSeqDao) GenerateNextSeq(convID int64) (int64, error) {
	key := rediskey.ConvSeq{ID: convID}

	luaScript := `
		if redis.call('exists', KEYS[1]) == 1 then
			return redis.call('incr', KEYS[1])
		else
			return -1
		end
	`

	// 1.从redis递增
	newSeq, err := dao.rdb.Eval(context.Background(), luaScript, []string{key.String()}).Int64()
	if err != nil {
		utils.LogWarn("从redis中找 conv:%d 的下一个消息seq失败，开始从数据库中找: error=%v", convID, err)
	}
	if newSeq != -1 {
		return newSeq, nil
	}

	var curSeq int64
	err = dao.db.Model(&model.IMMessage{}).
		Where("conv_id = ?", convID).
		Pluck("COALESCE(MAX(seq), 0)", &curSeq).Error
	if err != nil {
		return 0, err
	}

	newSeq = curSeq + 1

	err = dao.rdb.Set(context.Background(), key.String(), newSeq, key.TTL()).Err()
	if err != nil {
		return 0, err
	}

	return newSeq, nil
}
