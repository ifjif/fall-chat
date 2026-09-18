package data

import (
	"context"
	"errors"
	"fmt"
	"time"
	"zzc/fall-chat/src/database"
	"zzc/fall-chat/src/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ConversationDao struct {
	BaseDao
}

func NewConversationDao(db *gorm.DB) *ConversationDao {
	bd := BaseDao{db: db, rdb: database.Redis()}
	return &ConversationDao{
		BaseDao: bd,
	}
}

func (dao *ConversationDao) ListByUserID(userID int64) ([]model.IMConversationRow, error) {
	// 携带
	// peer_avatar
	// peer_nickname

	var rows []model.IMConversationRow
	err := dao.db.Table("im_conversation_members cm").
		Joins("inner join im_conversations cs on cm.conv_id = cs.id").
		Where("cm.user_id = ?", userID).
		Select("cs.*, cm.peer_id,cm.last_read_seq").
		Find(&rows).Error

	reqKeys := make([]string, 0, len(rows)*2)
	for i := range rows {
		ktype := "user"
		if rows[i].IsGroupConType() {
			ktype = "group"
		}
		avatarKey := fmt.Sprintf("im:%s:avatar:%d", ktype, rows[i].PeerID)
		nicknameKey := fmt.Sprintf("im:%s:nickname:%d", ktype, rows[i].PeerID)

		reqKeys = append(reqKeys, avatarKey, nicknameKey)
	}

	cmd := dao.rdb.MGet(context.Background(), reqKeys...)
	err = cmd.Err()
	if err != nil && err != redis.Nil {
		// TODO:
		// redis server error
	}

	result, _ := cmd.Result()
	var missingValue []int64
	missingMap := make(map[int64]int)

	for i := 0; i < len(result); i += 2 {
		avatarIndex := i
		nicknameIndex := i + 1

		var avatar, nickname string

		if result[avatarIndex] != nil {
			avatar, _ = result[avatarIndex].(string)
		}

		if result[nicknameIndex] != nil {
			nickname, _ = result[nicknameIndex].(string)
		}

		if avatar != "" && nickname != "" {
			rows[i/2].PeerAvatar = avatar
			rows[i/2].PeerNickname = nickname
		}
		peerID := rows[i/2].PeerID
		missingValue = append(missingValue, peerID)
		// 在rows中的索引
		missingMap[peerID] = i / 2
	}

	if len(missingValue) == 0 {
		return rows, nil
	}

	var infoResult []model.UserOrGroupBaseInfoRow
	err = dao.db.Raw("? union all ?",
		dao.db.Select("id,nickname as peer_nickname, avatar").Model(&model.User{}).Where("id in ?", missingValue),
		dao.db.Select("id, group_name as peer_nickname,avatar").Model(&model.GroupInfo{}).Where("id in ?", missingValue),
	).Scan(&infoResult).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 附加信息并设置缓存
	p := dao.rdb.Pipeline()
	duration := 1 * time.Hour
	ctx := context.Background()
	for i := range infoResult {
		peerID := infoResult[i].ID
		idx := missingMap[peerID]
		avatar := infoResult[i].Avatar
		nickname := infoResult[i].PeerNickname
		rows[idx].PeerAvatar = avatar
		rows[idx].PeerNickname = nickname

		p.SetEx(ctx, reqKeys[idx*2], avatar, duration)
		p.SetEx(ctx, reqKeys[idx*2+1], nickname, duration)
	}

	if len(infoResult) != 0 {
		p.Exec(ctx)
	}

	return rows, err
}

func (dao *ConversationDao) SingleByConvID(convID, userID int64) (*model.IMConversationRow, error) {
	// 携带
	// peer_avatar
	// peer_nickname

	var row model.IMConversationRow
	err := dao.db.Table("im_conversations cs").
		Joins("inner join im_conversation_members cm on cm.conv_id = cs.id and cm.user_id = ?", userID).
		Select("cs.*, cm.peer_id,cm.last_read_seq").
		Where("cs.id = ?", convID).
		Take(&row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	typeStr := "user"
	if row.IsGroupConType() {
		typeStr = "group"
	}

	avatarKey := fmt.Sprintf("im:%s:avatar:%d", typeStr, row.PeerID)
	nicknameKey := fmt.Sprintf("im:%s:nickname:%d", typeStr, row.PeerID)

	cmd := dao.rdb.MGet(context.Background(), avatarKey, nicknameKey)
	err = cmd.Err()
	if err != nil && err != redis.Nil {
		// TODO:
		fmt.Println("redis 服务错误", err)
	}

	result, _ := cmd.Result()
	if len(result) == 2 {
		if result[0] != nil && result[1] != nil {
			row.PeerAvatar = result[0].(string)
			row.PeerNickname = result[1].(string)
			return &row, nil
		}
	}

	var infoResult model.UserOrGroupBaseInfoRow
	if row.IsGroupConType() {
		err = dao.db.Select("id, group_name as peer_nickname,avatar").Model(&model.GroupInfo{}).
			Where("id = ?", row.PeerID).
			First(&infoResult).Error
	} else {
		err = dao.db.Select("id,nickname as peer_nickname, avatar").Model(&model.User{}).
			Where("id = ?", row.PeerID).
			First(&infoResult).Error
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	row.PeerAvatar = infoResult.Avatar
	row.PeerNickname = infoResult.PeerNickname

	// 存入redis
	p := dao.rdb.Pipeline()
	ctx := context.Background()
	duration := 1 * time.Hour
	p.SetEx(ctx, avatarKey, infoResult.Avatar, duration)
	p.SetEx(ctx, nicknameKey, infoResult.PeerNickname, duration)
	p.Exec(ctx)

	return &row, err
}
