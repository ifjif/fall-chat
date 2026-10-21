package converter

import (
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToUserOrGroupQuery(req *request.UserOrGroup) *query.UserOrGroup {
	return &query.UserOrGroup{
		ID:         req.ID,
		Username:   req.Username,
		Name:       req.Name,
		NextCursor: req.NextCuror,
	}
}

func ToUserOrGroupDataResult(row *model.UserOrGroupRelationRow) *response.UserOrGroupDataResult {
	return &response.UserOrGroupDataResult{
		Type:           row.RelationType,
		ID:             row.ID,
		Name:           row.Name,
		Avatar:         row.Avatar,
		Nickname:       row.Nickname,
		Description:    row.Description,
		JoinMode:       row.JoinMode,
		RelationStatus: row.GetRelationStatusText(),
	}
}
