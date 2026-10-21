package cache

import (
	"context"
	"strconv"
	"time"
	"zzc/fall-chat/src/keys/ctxkey"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/utils"

	"github.com/redis/go-redis/v9"
)

type GroupMemberCache interface {
	FindUserIDs(ctx context.Context, groupID int64) ([]int64, error)
	SetUserIDs(ctx context.Context, groupID int64, userIDs []int64)
	SetUserID(ctx context.Context, groupID, userID int64)
	DelUserID(ctx context.Context, groupID, userID int64)
}

type groupMemberCache struct {
	rdb *redis.Client
}

func NewGroupMemberCache(rdb *redis.Client) GroupMemberCache {
	return &groupMemberCache{
		rdb: rdb,
	}
}

func (c *groupMemberCache) FindUserIDs(ctx context.Context, groupID int64) ([]int64, error) {
	key := rediskey.GroupUserIDs{GroupID: groupID}

	res, err := c.rdb.SMembers(context.Background(), key.String()).Result()
	if err != nil {
		return nil, err
	}

	var result []int64

	for _, v := range res {
		iv, _ := strconv.ParseInt(v, 0, 64)
		result = append(result, iv)
	}
	return result, nil
}

func (c *groupMemberCache) SetUserIDs(ctx context.Context, groupID int64, userIDs []int64) {
	logger := ctxkey.GetLogger(ctx)
	key := rediskey.GroupUserIDs{GroupID: groupID}
	luascript := `
		redis.call('sadd', KEYS[1], unpack(ARGV, 2, #ARGV))
		return redis.call('expire', KEYS[1], tonumber(ARGV[1]))
	`

	args := []any{
		key.TTL() / time.Second,
	}
	args = append(args, utils.ToAnySlice(userIDs)...)
	err := c.rdb.Eval(context.Background(), luascript, []string{key.String()}, args...).Err()
	if err != nil {
		logger.Error("[redis]", "reason", "设置group的所有用户id到缓存失败", "group_id", groupID, "error", err)
	}
}

func (c *groupMemberCache) SetUserID(ctx context.Context, groupID, userID int64) {
	logger := ctxkey.GetLogger(ctx)
	key := rediskey.GroupUserIDs{GroupID: groupID}
	err := c.rdb.SAdd(ctx, key.String(), userID).Err()
	if err != nil {
		logger.Error("[redis]",
			"reason", "设置group的一个用户id到缓存失败",
			"group_id", groupID,
			"user_id", userID,
			"error", err)
	}
}

func (c *groupMemberCache) DelUserID(ctx context.Context, groupID, userID int64) {
	logger := ctxkey.GetLogger(ctx)
	key := rediskey.GroupUserIDs{GroupID: groupID}
	err := c.rdb.SRem(ctx, key.String(), userID).Err()
	if err != nil {
		logger.Error("[redis]",
			"reason", "删除group的一个用户id失败",
			"group_id", groupID,
			"user_id", userID,
			"error", err)
	}
}
