package cache

import (
	"context"
	"fmt"
	"time"
	"zzc/fall-chat/src/keys/ctxkey"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/model/cached"
	"zzc/fall-chat/src/types"

	"github.com/redis/go-redis/v9"
)

type RTCCache interface {
	CreateRoom(ctx context.Context, creatorID int64, typ cached.RtcRoomType, maxMember int) (string, error)
	GetRoom(ctx context.Context, roomID string) (*cached.RtcRoom, error)
	UpdateRoomState(ctx context.Context, roomID string, state cached.RtcRoomState) error
	GetRoomState(ctx context.Context, roomID string) string
	SetPeerRoom(ctx context.Context, userID int64, roomID string) error
	GetPeerRoom(ctx context.Context, userID int64) (string, error)
	DelPeerRoom(ctx context.Context, userID int64)
}

type rtcCache struct {
	rdb *redis.Client
}

func NewRTCCache(rdb *redis.Client) RTCCache {
	return &rtcCache{
		rdb: rdb,
	}
}

func (c *rtcCache) CreateRoom(ctx context.Context, creatorID int64, typ cached.RtcRoomType, maxMember int) (string, error) {
	now := time.Now()
	roomID := fmt.Sprintf("room_%d_%d", creatorID, now.Nanosecond())
	key := rediskey.RtcRoom{
		RoomID: roomID,
	}

	room := cached.RtcRoom{
		RoomID:      roomID,
		CreatorID:   creatorID,
		Type:        typ,
		MaxMember:   maxMember,
		CreatedTime: types.RedisTime{Time: now},
		ExpiredTime: types.RedisTime{Time: now.Add(key.TTL())},
		State:       cached.RtcRoomStatePending,
	}

	pipe := c.rdb.TxPipeline()

	pipe.HSet(ctx, key.String(), room)
	pipe.Expire(ctx, key.String(), key.TTL())

	if _, err := pipe.Exec(ctx); err != nil {
		return "", err
	}

	return roomID, nil
}

func (c *rtcCache) GetRoom(ctx context.Context, roomID string) (*cached.RtcRoom, error) {
	key := rediskey.RtcRoom{RoomID: roomID}
	cmd := c.rdb.HGetAll(ctx, key.String())
	var v cached.RtcRoom
	err := cmd.Scan(&v)

	return &v, err
}

func (c *rtcCache) UpdateRoomState(ctx context.Context, roomID string, state cached.RtcRoomState) error {
	key := rediskey.RtcRoom{RoomID: roomID}
	pipe := c.rdb.TxPipeline()

	pipe.HSet(ctx, key.String(), "state", string(state)).Err()
	pipe.Expire(ctx, key.String(), key.TTL())

	_, err := pipe.Exec(ctx)
	return err
}

func (c *rtcCache) GetRoomState(ctx context.Context, roomID string) string {
	logger := ctxkey.GetLogger(ctx)
	key := rediskey.RtcRoom{RoomID: roomID}
	state, err := c.rdb.HGet(ctx, key.String(), "state").Result()
	if err != nil && err != redis.Nil {
		logger.Error("[redis]", "reason", "找rtc_room的state", "key", key.String(), "error", err)
	}

	return state
}

func (c *rtcCache) SetPeerRoom(ctx context.Context, userID int64, roomID string) error {
	key := rediskey.RtcPeerRoom{UserID: userID}
	return c.rdb.Set(ctx, key.String(), roomID, key.TTL()).Err()
}

func (c *rtcCache) GetPeerRoom(ctx context.Context, userID int64) (string, error) {
	key := rediskey.RtcPeerRoom{UserID: userID}
	result, err := c.rdb.Get(ctx, key.String()).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	return result, nil
}

func (c *rtcCache) DelPeerRoom(ctx context.Context, userID int64) {
	logger := ctxkey.GetLogger(ctx)
	key := rediskey.RtcPeerRoom{UserID: userID}
	err := c.rdb.Del(ctx, key.String()).Err()
	if err != nil {
		logger.Error("[redis]", "reason", "从缓存中删除peer_room", "key", key.String(), "error", err)
	}
}
