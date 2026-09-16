package service

import (
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

type ConversationMemberService struct {
	BaseService
	dao *data.ConversationMemberDao
}

func NewConversationMemberService(name string) *ConversationMemberService {
	bs := BaseService{name: name}

	return &ConversationMemberService{
		BaseService: bs,
	}
}

func (cms *ConversationMemberService) Inject() {
	cms.BaseService.Inject()
	cms.dao = data.GetConversationMemberDao(cms.db)
}

func (cms *ConversationMemberService) ListBrief4Conv(qu *query.ConvMembListBrief4ConvQuery) (map[string]*model.IMConvMembBriefRow, error) {
	convID, err := cms.sf.Decode(qu.ConvID)
	if err != nil {
		return nil, cms.Wrap("解码 conversation id", err, "id", qu.ID)
	}

	rows, err := cms.dao.ListBriefByConvID(convID)
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

func (cms *ConversationMemberService) ListBriefBatch(qu *query.ConvMembListBriefBatchQuery) (map[string]*model.IMConvMembBriefRow, error) {
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
	rows, err := cms.dao.ListBriefBatchByConvIDAndMemberIDs(convID, memberIDs, memberPIDs)
	if err != nil {
		return nil, cms.Wrap("找会话成员基本信息", err, "id", qu.ID,
			"conv_id", qu.ConvID, "member_pids", memberPIDs)
	}

	result := make(map[string]*model.IMConvMembBriefRow, len(rows))
	for i := range rows {
		row := rows[i]
		result[row.PublicID] = &row
	}

	return result, nil
}
