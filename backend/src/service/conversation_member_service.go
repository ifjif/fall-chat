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
		return nil, cms.LogWarn("ListBrief4Conv", err, "id", qu.ID, "reason", "解码 conversation id 失败")
	}

	rows, err := cms.dao.ListBriefByConvID(convID)
	if err != nil {
		return nil, cms.LogWarn("ListBriefByConvID", err, "id", qu.ID, "reason", "找会话成员基本信息失败")
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
		return nil, cms.LogWarn("ListBriefByConvID", err, "id", qu.ID,
			"reason", "解码会话id失败")
	}
	memberPIDs := qu.MemberIDs
	memberIDs := make([]int64, 0, len(memberPIDs))

	for _, pID := range memberPIDs {
		id, err := cms.sf.Decode(pID)
		if err != nil {
			return nil, cms.LogWarn("ListBriefBatch", err, "id", id, "memberPIDs", memberIDs,
				"reason", "解析会话成员id失败")
		}

		memberIDs = append(memberIDs, id)
	}
	rows, err := cms.dao.ListBriefBatchByConvIDAndMemberIDs(convID, memberIDs, memberPIDs)
	if err != nil {
		return nil, cms.LogWarn("ListBriefBatch", err, "id", qu.ID,
			"conv_id", qu.ConvID, "member_pids", memberPIDs, "reason", "找会话成员基本信息失败")
	}

	result := make(map[string]*model.IMConvMembBriefRow, len(rows))
	for i := range rows {
		row := rows[i]
		result[row.PublicID] = &row
	}

	return result, nil
}
