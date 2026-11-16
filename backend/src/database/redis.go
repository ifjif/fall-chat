package database

import (
	"log"
	"strconv"
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
	addr, addrExists := utils.GetEnv("REDIS_ADDR")
	pass, passExists := utils.GetEnv("REDIS_PASS")
	if !addrExists || !passExists {
		utils.LogError("redis addr or password 不存在，无法使用redis")
		log.Fatal()
	}
	dbStr, _ := utils.GetEnv("REDIS_DB")
	protocolStr, _ := utils.GetEnv("REDIS_PROTOCOL")

	db, err := strconv.ParseInt(dbStr, 10, 0)
	if err != nil {
		utils.LogWarn("redis db 解析错误，采用默认0")
		db = 0
	}
	protocol, err := strconv.ParseInt(protocolStr, 10, 0)
	if err != nil {
		utils.LogWarn("redis protocol 解析错误，采用默认2")
		protocol = 2
	}

	return &redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       int(db),
		Protocol: int(protocol),
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
