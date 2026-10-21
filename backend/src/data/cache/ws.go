package cache

import (
	"context"
	"zzc/fall-chat/src/host"
	"zzc/fall-chat/src/keys/ctxkey"
	rediskey "zzc/fall-chat/src/keys/redis_key"

	"github.com/redis/go-redis/v9"
)

type WSCache interface {
	AddGrpcAddr(ctx context.Context, id string)
	DelGrpcAddr(ctx context.Context, id string)
	GetGrpcAddr(ctx context.Context, id string) (addr string, err error)
}

type wsCache struct {
	rdb *redis.Client
}

func NewWSCache(rdb *redis.Client) WSCache {
	return &wsCache{rdb: rdb}
}

func (c *wsCache) AddGrpcAddr(ctx context.Context, id string) {
	logger := ctxkey.GetLogger(ctx)
	key := rediskey.GrpcAddr{ID: id}
	grpcAddr := host.GRPCAddr
	// 5min
	err := c.rdb.Set(ctx, key.String(), grpcAddr, key.TTL()).Err()
	if err != nil {
		logger.Error("[redis]", "reason", "设置grpc address失败", "error", err, "key", key)
	}
}

func (c *wsCache) DelGrpcAddr(ctx context.Context, id string) {
	key := rediskey.GrpcAddr{ID: id}
	c.rdb.Del(ctx, key.String())
}

func (c *wsCache) GetGrpcAddr(ctx context.Context, id string) (addr string, err error) {
	key := rediskey.GrpcAddr{ID: id}

	cmd := c.rdb.Get(ctx, key.String())

	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		return "", cmd.Err()
	}
	addr, _ = cmd.Result()

	return addr, nil
}
