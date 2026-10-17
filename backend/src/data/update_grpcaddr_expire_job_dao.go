package data

import (
	"context"
	"sync"
	"zzc/fall-chat/src/database"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/utils"
)

type UpdateGrpcaddrExpireJobDao struct {
	BaseDao
}

var (
	updateGrpcaddrExpireJobDaoOnce sync.Once
	updateGrpcaddrExpireJobDao     *UpdateGrpcaddrExpireJobDao
)

func GetUpdateGrpcaddrExpireJobDao() *UpdateGrpcaddrExpireJobDao {
	updateGrpcaddrExpireJobDaoOnce.Do(func() {
		updateGrpcaddrExpireJobDao = &UpdateGrpcaddrExpireJobDao{}
		updateGrpcaddrExpireJobDao.rdb = database.Redis()
	})

	return updateGrpcaddrExpireJobDao
}

func (dao *UpdateGrpcaddrExpireJobDao) UpdateGrpcAddrExpire2Redis(id string) {
	key := rediskey.GrpcAddr{ID: id}
	ok, err := dao.rdb.Expire(context.Background(), key.String(), key.TTL()).Result()
	if err != nil {
		utils.LogError("设置%s 失败：error=%v", key, err)
		return
	}

	if !ok {
		utils.LogError("键 %s 不存在或已过期", key)
	}
}
