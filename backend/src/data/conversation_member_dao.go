package data

import (
	"context"
	"fmt"
	"sync"
	"time"
	"zzc/fall-chat/src/database"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/utils"

	"gorm.io/gorm"
)

type ConversationMemberDao struct {
	BaseDao
}

var (
	conversationMemberDaoOnce sync.Once
	conversationMemberDao     *ConversationMemberDao
)

func GetConversationMemberDao(db *gorm.DB) *ConversationMemberDao {
	conversationMemberDaoOnce.Do(func() {
		bd := BaseDao{db: db, rdb: database.Redis()}
		conversationMemberDao = &ConversationMemberDao{
			BaseDao: bd,
		}
	})

	return conversationMemberDao
}

func (dao *ConversationMemberDao) listBrief(convID int64, memberIDs []int64) ([]model.IMConvMembBriefRow, error) {
	var rows []model.IMConvMembBriefRow
	db := dao.db.Table("im_conversation_members icm").
		Joins("inner join users u on icm.user_id = u.id").
		Joins("left join group_members gm on gm.group_id = icm.peer_id and gm.user_id = icm.user_id").
		Where("icm.conv_id = ?", convID)

	if len(memberIDs) != 0 {
		db.Where("icm.user_id in ?", memberIDs)
	}

	err := db.Select("u.id, u.public_id, u.nickname, gm.nickname as in_group_nickname, u.avatar").
		Find(&rows).Error

	return rows, err
}

func (dao *ConversationMemberDao) ListBriefByConvID(convID int64) ([]model.IMConvMembBriefRow, error) {
	rows, err := dao.listBrief(convID, nil)
	return rows, err
}

func (dao *ConversationMemberDao) ListBriefBatchByConvIDAndMemberIDs(convID int64, memberIDs []int64, memberPIDs []string) ([]model.IMConvMembBriefRow, error) {
	// 从redis中找 批次找
	// 剩下没有的从数据库找
	// 数据库中找到的回写到redis
	// 最后返回

	// im:mem:nickname:${memberID}
	// im:mem:avatar:${memberID}
	// im:mem:ingroupname:${memberID}
	if len(memberIDs) == 0 {
		return nil, nil
	}
	nicknameKeyPrefix := "im:mem:nickname"
	avatarKeyPrefix := "im:mem:avatar"
	ingroupnameKeyPrefix := "im:mem:ingroupname"

	keys := make([]string, 0, len(memberIDs)*3)

	for _, id := range memberIDs {
		keys = append(keys, fmt.Sprintf("%s:%d", nicknameKeyPrefix, id))
		keys = append(keys, fmt.Sprintf("%s:%d", avatarKeyPrefix, id))
		keys = append(keys, fmt.Sprintf("%s:%d:%d", ingroupnameKeyPrefix, convID, id))
	}

	mgetCmd := dao.rdb.MGet(context.Background(), keys...)

	if err := mgetCmd.Err(); err != nil {
		return nil, err
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

	if len(missing) == 0 {
		return result, nil
	}

	rows, err := dao.listBrief(convID, missing)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return result, nil
	}

	result = append(result, rows...)

	expiration := 1 * time.Hour

	missKeyVals := make([]string, 0, len(rows)*6)
	for i := range rows {
		id := rows[i].ID
		missKeyVals = append(missKeyVals, fmt.Sprintf("%s:%d", nicknameKeyPrefix, id))
		missKeyVals = append(missKeyVals, rows[i].Nickname)

		missKeyVals = append(missKeyVals, fmt.Sprintf("%s:%d", avatarKeyPrefix, id))
		missKeyVals = append(missKeyVals, rows[i].Avatar)

		missKeyVals = append(missKeyVals, fmt.Sprintf("%s:%d:%d", ingroupnameKeyPrefix, convID, id))
		missKeyVals = append(missKeyVals, rows[i].InGropNickname)
	}

	pipe := dao.rdb.Pipeline()
	pipe.MSet(context.Background(), missKeyVals)
	for i := 0; i < len(missKeyVals); i += 2 {
		key := missKeyVals[i]
		pipe.Expire(context.Background(), key, expiration)
	}

	_, err = pipe.Exec(context.Background())
	if err != nil {
		// 记录缓存设置失败
		utils.LogWarn("ConversationMemberDao[ListBriefBatchByConvIDAndMemberIDs] 缓存设置失败 | context [%v]", missing)
	}

	return result, nil
}
