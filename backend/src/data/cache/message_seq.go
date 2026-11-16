package cache

import (
	"context"
	"zzc/fall-chat/src/keys/ctxkey"
	rediskey "zzc/fall-chat/src/keys/redis_key"

	"github.com/redis/go-redis/v9"
)

type MessageSeqCache interface {
	IncrConvSeq(ctx context.Context, convID int64) (seq int64, ok bool)
	SetConvSeq(ctx context.Context, convID, seq int64) error
}

type messageSeqCache struct {
	rdb *redis.Client
}

func NewMessageSeqCache(rdb *redis.Client) MessageSeqCache {
	return &messageSeqCache{
		rdb: rdb,
	}
}

func (c *messageSeqCache) IncrConvSeq(ctx context.Context, convID int64) (seq int64, ok bool) {
	logger := ctxkey.GetLogger(ctx)
	key := rediskey.ConvSeq{ID: convID}

	luaScript := `
		if redis.call('exists', KEYS[1]) == 1 then
			return redis.call('incr', KEYS[1])
		else
			return -1
		end
	`

	// 从redis递增
	newSeq, err := c.rdb.Eval(context.Background(), luaScript, []string{key.String()}).Int64()
	if err != nil {
		logger.Error("[redis]", "reason", "自增会话消息seq失败", "convId", convID, "error", err)
		return -1, false
	}

	if newSeq == -1 {
		return -1, false
	}

	return newSeq, true
}

func (c *messageSeqCache) SetConvSeq(ctx context.Context, convID, seq int64) error {
	key := rediskey.ConvSeq{ID: convID}
	return c.rdb.Set(context.Background(), key.String(), seq, key.TTL()).Err()
}
