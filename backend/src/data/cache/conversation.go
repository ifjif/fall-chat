package cache

import (
	"context"
	"encoding/json"
	"time"
	"zzc/fall-chat/src/keys/ctxkey"
	rediskey "zzc/fall-chat/src/keys/redis_key"
	"zzc/fall-chat/src/model"

	"github.com/redis/go-redis/v9"
)

type ConversationCache interface {
	// 找会话的简单信息
	// id public_id conv_type
	GetBrief(ctx context.Context, singleKey string) (*model.IMConversation, error)
	// 设置简单信息
	// id public_id conv_type
	// conv中需要存在 SingleKey 来当缓存key
	SetBrief(ctx context.Context, conv model.IMConversation)

	// 附加会话 nickname,avatar
	// 返回未附加的 会话id-> rows中的idx
	AppendInfo(ctx context.Context, rows []*model.IMConversationRow) (miss map[int64]int, err error)

	// 设置会话的 nickname,avatar信息
	SetInfo(ctx context.Context, infoResult []model.UserOrGroupBaseInfoRow) error
}

type conversationCache struct {
	rdb *redis.Client
}

func NewConversationCache(rdb *redis.Client) ConversationCache {
	return &conversationCache{
		rdb: rdb,
	}
}

func (c *conversationCache) GetBrief(ctx context.Context, singleKey string) (*model.IMConversation, error) {
	key := rediskey.Conv{SingleKey: singleKey}

	res, err := c.rdb.Get(context.Background(), key.String()).Result()
	if err != nil && err != redis.Nil {
		// 记录错误
		// logger.Error("[redis]", "获取会话id,public_id,conv_type失败")
		return nil, err
	}
	if err == redis.Nil {
		return nil, nil
	}

	var conv model.IMConversation
	err = json.Unmarshal([]byte(res), &conv)
	if err != nil {
		return nil, err
	}
	// utils.LogError("json反序列化错误：%v", err)
	return &conv, nil
}

func (c *conversationCache) SetBrief(ctx context.Context, conv model.IMConversation) {
	logger := ctxkey.GetLogger(ctx)

	key := rediskey.Conv{SingleKey: conv.SingleKey}
	cacheInfo := map[string]any{
		"id":        conv.ID,
		"public_id": conv.PublicID,
		"conv_type": conv.ConvType,
	}
	value, err := json.Marshal(cacheInfo)
	if err != nil {
		logger.Error("[json]", "reason", "序列化失败", "error", err)
		return
	}

	err = c.rdb.Set(context.Background(), key.String(), value, key.TTL()).Err()
	if err != nil {
		logger.Error("[redis]", "reason", "设置会话id,public_id,conv_type到缓存失败")
	}
}

func (c *conversationCache) AppendInfo(ctx context.Context, rows []*model.IMConversationRow) (miss map[int64]int, err error) {
	reqKeys := make([]string, 0, len(rows)*2)

	userAvatarKey := rediskey.ImUserAvatar{}
	userNicknameKey := rediskey.ImUserNickname{}
	groupAvatarKey := rediskey.ImGroupAvatar{}
	groupNicknameKey := rediskey.ImGroupNickname{}
	for i := range rows {
		userAvatarKey.ID = rows[i].PeerID
		userNicknameKey.ID = rows[i].PeerID
		groupAvatarKey.ID = rows[i].PeerID
		groupNicknameKey.ID = rows[i].PeerID

		avatarKey := userAvatarKey.String()
		nicknameKey := userNicknameKey.String()
		if rows[i].IsGroupConType() {
			avatarKey = groupAvatarKey.String()
			nicknameKey = groupNicknameKey.String()
		}

		reqKeys = append(reqKeys, avatarKey, nicknameKey)
	}

	missingMap := make(map[int64]int)
	res, err := c.rdb.MGet(context.Background(), reqKeys...).Result()
	if err != nil {
		// redis server error
		return missingMap, err
	}

	for i := 0; i < len(res); i += 2 {
		avatarIndex := i
		nicknameIndex := i + 1

		var avatar, nickname string

		if res[avatarIndex] != nil {
			avatar, _ = res[avatarIndex].(string)
		}

		if res[nicknameIndex] != nil {
			nickname, _ = res[nicknameIndex].(string)
		}

		if avatar != "" && nickname != "" {
			rows[i/2].PeerAvatar = avatar
			rows[i/2].PeerNickname = nickname
		}
		peerID := rows[i/2].PeerID

		// 在rows中的索引
		missingMap[peerID] = i / 2
	}

	return missingMap, nil
}

func (c *conversationCache) SetInfo(ctx context.Context, infoResult []model.UserOrGroupBaseInfoRow) error {
	userAvatarKey := rediskey.ImUserAvatar{}
	userNicknameKey := rediskey.ImUserNickname{}
	groupAvatarKey := rediskey.ImGroupAvatar{}
	groupNicknameKey := rediskey.ImGroupNickname{}

	// 附加信息并设置缓存
	p := c.rdb.Pipeline()
	avatarDuration := time.Hour
	nicknameDuration := time.Hour

	for i := range infoResult {
		peerID := infoResult[i].ID
		userAvatarKey.ID = peerID
		userNicknameKey.ID = peerID
		groupAvatarKey.ID = peerID
		groupNicknameKey.ID = peerID

		var avatarKey string
		var nicknameKey string

		if infoResult[i].IsGroup() {
			avatarKey = groupAvatarKey.String()
			nicknameKey = groupNicknameKey.String()
			avatarDuration = groupAvatarKey.TTL()
			nicknameDuration = groupNicknameKey.TTL()
		} else if infoResult[i].IsUser() {
			avatarKey = userAvatarKey.String()
			nicknameKey = userNicknameKey.String()
			avatarDuration = userAvatarKey.TTL()
			nicknameDuration = userNicknameKey.TTL()
		}

		avatar := infoResult[i].Avatar
		nickname := infoResult[i].PeerNickname

		if avatarKey != "" {
			p.SetEx(ctx, avatarKey, avatar, avatarDuration)
		}
		if nicknameKey != "" {
			p.SetEx(ctx, nicknameKey, nickname, nicknameDuration)
		}
	}

	if len(infoResult) != 0 {
		_, err := p.Exec(ctx)
		return err
	}

	return nil
}
