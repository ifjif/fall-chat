package service

import (
	"context"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/keys/ctxkey"
)

type GroupMember struct {
	BaseService
	store store.GroupMemberStore
}

func NewGroupMember(name string) *GroupMember {
	bs := BaseService{name: name}
	return &GroupMember{
		BaseService: bs,
	}
}

func (gm *GroupMember) Inject() {
	gm.store = data.GetBaseDao()
}

func (gm *GroupMember) FindUserIDs(ctx context.Context, groupID int64) ([]int64, error) {
	logger := ctxkey.GetLogger(ctx)
	// 缓存
	ids, err := gm.store.
		GroupMemberCache().
		FindUserIDs(ctx, groupID)
	if err != nil {
		logger.Error("[redis]", "找群的所有用户ids失败", "group_id", groupID, "error", err)
	}

	if len(ids) > 0 {
		return ids, nil
	}

	// 从数据库
	ids, err = gm.store.
		GroupMemberRepos().
		GroupMember.
		FindUserIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// 设置缓存
	gm.store.
		GroupMemberCache().
		SetUserIDs(ctx, groupID, ids)

	return ids, nil
}
