package cache

import (
	"context"
	"zzc/fall-chat/src/keys/ctxkey"
	rediskey "zzc/fall-chat/src/keys/redis_key"

	"github.com/redis/go-redis/v9"
)

type EmojiCache interface {
	SetPublishItemCount(ctx context.Context, userID, packID uint64, count int) error
	GetPublishItemCount(ctx context.Context, userID, packID uint64) int
}

type emojiCache struct {
	rdb *redis.Client
}

func NewEmojiCache(rdb *redis.Client) EmojiCache {
	return &emojiCache{
		rdb: rdb,
	}
}

func (c *emojiCache) SetPublishItemCount(ctx context.Context, userID, packID uint64, count int) error {
	key := rediskey.PublishEmojiItemCount{
		UserID: userID,
		PackID: packID,
	}

	return c.rdb.SetNX(ctx, key.String(), count, key.TTL()).Err()
}

func (c *emojiCache) GetPublishItemCount(ctx context.Context, userID, packID uint64) int {
	logger := ctxkey.GetLogger(ctx)
	key := rediskey.PublishEmojiItemCount{
		UserID: userID,
		PackID: packID,
	}

	count, err := c.rdb.Get(ctx, key.String()).Int()
	if err != nil && err != redis.Nil {
		logger.Error("[redis]", "reason", "从缓存中找emoji_item推送数量", "error", err,
			"user_id", userID, "pack_id", packID)
	}

	return count
}
