package service

import (
	"context"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

type ConversationMember struct {
	BaseService
	store store.ConvMembStore
}

func NewConversationMember(name string) *ConversationMember {
	bs := BaseService{name: name}

	return &ConversationMember{
		BaseService: bs,
	}
}

func (cms *ConversationMember) Inject() {
	cms.BaseService.Inject()
	cms.store = data.GetBaseDao()
}

func (cms *ConversationMember) ListBrief4Conv(ctx context.Context, qu *query.ConvMembListBrief4Conv) (map[string]*model.IMConvMembBriefRow, error) {
	convID, err := cms.sf.Decode(qu.ConvID)
	if err != nil {
		return nil, cms.Wrap("解码 conversation id", err, "id", qu.ID)
	}

	rows, err := cms.store.
		ConvMembRepos().
		ConvMemb.
		ListBriefByConvID(ctx, convID)
	if err != nil {
		return nil, cms.Wrap("找会话成员基本信息", err, "id", qu.ID)
	}

	result := make(map[string]*model.IMConvMembBriefRow, len(rows))
	for i := range rows {
		row := rows[i]
		result[row.PublicID] = &row
	}

	return result, nil
}

func (cms *ConversationMember) ListBriefBatch(ctx context.Context, qu *query.ConvMembListBriefBatch) (map[string]*model.IMConvMembBriefRow, error) {
	// 从redis中找 批次找
	// 剩下没有的从数据库找
	// 数据库中找到的回写到redis
	// 最后返回
	convID, err := cms.sf.Decode(qu.ConvID)
	if err != nil {
		return nil, cms.Wrap("解码会话id", err, "id", qu.ID)
	}
	memberPIDs := qu.MemberIDs
	memberIDs := make([]int64, 0, len(memberPIDs))

	for _, pID := range memberPIDs {
		id, err := cms.sf.Decode(pID)
		if err != nil {
			return nil, cms.Wrap("解析会话成员id", err, "id", id, "memberPIDs", memberIDs)
		}

		memberIDs = append(memberIDs, id)
	}

	// 从缓存中找
	rows, miss, err := cms.store.
		ConvMembCache().
		FindBriefInfo(ctx, convID, memberIDs, memberPIDs)
	if err != nil {
		// TODO:
	}

	var missPIDs []string
	for _, id := range miss {
		missPIDs = append(missPIDs, cms.sf.Encode(id))
	}

	// 找忽略的
	missRows, err := cms.store.
		ConvMembRepos().
		ConvMemb.
		ListBriefBatchByConvIDAndMemberIDs(ctx, convID, miss)
	if err != nil {
		return nil, cms.Wrap("找会话成员基本信息", err, "id", qu.ID,
			"conv_id", qu.ConvID, "member_pids", memberPIDs)
	}

	// 存入缓存
	err = cms.store.
		ConvMembCache().
		SetBriefInfo(ctx, convID, missRows)
	if err != nil {
		// TODO:
	}

	rows = append(rows, missRows...)

	result := make(map[string]*model.IMConvMembBriefRow, len(rows))
	for i := range rows {
		row := rows[i]
		result[row.PublicID] = &row
	}

	return result, nil
}
