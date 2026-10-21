package cache

import (
	"context"
	"time"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/model"

	"github.com/redis/go-redis/v9"
)

type ConvMembCache interface {
	FindBriefInfo(ctx context.Context, convID int64, memberIDs []int64, memberPIDs []string) ([]model.IMConvMembBriefRow, []int64, error)

	SetBriefInfo(ctx context.Context, convID int64, rows []model.IMConvMembBriefRow) error
}

type convMembCache struct {
	rdb *redis.Client
}

func NewConvMembCache(rdb *redis.Client) ConvMembCache {
	return &convMembCache{rdb: rdb}
}

func (c *convMembCache) FindBriefInfo(ctx context.Context, convID int64, memberIDs []int64, memberPIDs []string) ([]model.IMConvMembBriefRow, []int64, error) {
	nicknameKey := rediskey.ImUserNickname{}
	avatarKey := rediskey.ImUserAvatar{}
	ingroupnameKey := rediskey.ImUserIngroupname{ConvID: convID}

	keys := make([]string, 0, len(memberIDs)*3)

	for _, id := range memberIDs {
		nicknameKey.ID = id
		avatarKey.ID = id
		ingroupnameKey.ID = id
		keys = append(keys, nicknameKey.String())
		keys = append(keys, avatarKey.String())
		keys = append(keys, ingroupnameKey.String())
	}

	mgetCmd := c.rdb.MGet(context.Background(), keys...)

	if err := mgetCmd.Err(); err != nil {
		return nil, nil, err
	}

	vals := mgetCmd.Val()

	var result []model.IMConvMembBriefRow
	var missing []int64

	for i, id := range memberIDs {
		nicknameIdx := i * 3
		avatarIdx := i*3 + 1
		ingroupnameIdx := i*3 + 2

		var nickname, avatar, ingroupname string

		if vals[nicknameIdx] != nil {
			nickname, _ = vals[nicknameIdx].(string)
		}
		if vals[avatarIdx] != nil {
			avatar, _ = vals[avatarIdx].(string)
		}
		if vals[ingroupnameIdx] != nil {
			ingroupname, _ = vals[ingroupnameIdx].(string)
		}

		if avatar == "" || (nickname == "" && ingroupname == "") {
			missing = append(missing, id)
			break
		}

		result = append(result, model.IMConvMembBriefRow{
			ID:             id,
			PublicID:       memberPIDs[i],
			Nickname:       nickname,
			InGropNickname: ingroupname,
			Avatar:         avatar,
		})
	}

	return result, missing, nil
}

func (c *convMembCache) SetBriefInfo(ctx context.Context, convID int64, rows []model.IMConvMembBriefRow) error {
	nicknameKey := rediskey.ImUserNickname{}
	avatarKey := rediskey.ImUserAvatar{}
	ingroupnameKey := rediskey.ImUserIngroupname{ConvID: convID}

	missKeyVals := make([]string, 0, len(rows)*6)
	missKeyExpires := []time.Duration{
		nicknameKey.TTL(),
		avatarKey.TTL(),
		ingroupnameKey.TTL(),
	}
	for i := range rows {
		id := rows[i].ID
		nicknameKey.ID = id
		avatarKey.ID = id
		ingroupnameKey.ID = id

		missKeyVals = append(missKeyVals, nicknameKey.String())
		missKeyVals = append(missKeyVals, rows[i].Nickname)

		missKeyVals = append(missKeyVals, avatarKey.String())
		missKeyVals = append(missKeyVals, rows[i].Avatar)

		missKeyVals = append(missKeyVals, ingroupnameKey.String())
		missKeyVals = append(missKeyVals, rows[i].InGropNickname)
	}

	pipe := c.rdb.Pipeline()
	pipe.MSet(context.Background(), missKeyVals)

	for i := 0; i < len(missKeyVals); i += 2 {
		key := missKeyVals[i]
		expiration := missKeyExpires[i%6/2]
		pipe.Expire(context.Background(), key, expiration)
	}

	_, err := pipe.Exec(context.Background())
	return err
}
