package database

import (
	"zzc/fall-chat/src/utils"

	"github.com/redis/go-redis/v9"
)

var redisInstance *redis.Client

func Redis() *redis.Client {
	return redisInstance
}

func initRedis() {
	redisOptions := getRedisOptions()
	redisInstance = redis.NewClient(redisOptions)

	utils.LogInfo("Redis初始成功:[Address:%s]", redisOptions.Addr)
}

func getRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     "redis7:6379",
		Password: "12345678",
		DB:       0,
		Protocol: 2,
	}
}

/*
	go func() {
		rdb := database.Redis()
		utils.LogDebug("开始发送redis消息")
		rdb.Set(context.Background(), "name", []byte("zhangsan"), 0)
		timer := time.NewTimer(2 * time.Second)
		<-timer.C
		sc := rdb.Get(context.Background(), "name")
		v := sc.Val()
		utils.LogDebug("从Redis中找到[name]=%s", v)
	}()


*/
