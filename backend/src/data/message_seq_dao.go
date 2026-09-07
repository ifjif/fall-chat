package data

import (
	"context"
	"fmt"
	"zzc/fall-chat/src/database"
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
	redisKey := fmt.Sprintf("conv:seq:%d", convID)

	// 1.从redis递增
	newSeq, err := dao.rdb.Incr(context.Background(), redisKey).Result()
	if err == nil {
		return newSeq, nil
	}

	utils.LogWarn("从redis中找 conv:%d 的下一个消息seq失败，开始从数据库中找", convID)

	var curSeq int64
	err = dao.db.Model(&model.IMMessage{}).
		Where("conv_id = ?", convID).
		Pluck("COALESCE(MAX(seq), 0)", &curSeq).Error
	if err != nil {
		return 0, err
	}

	newSeq = curSeq + 1

	err = dao.rdb.Set(context.Background(), redisKey, newSeq, 0).Err()
	if err != nil {
		return 0, err
	}

	return newSeq, nil
}
