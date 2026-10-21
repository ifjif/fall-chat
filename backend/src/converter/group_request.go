package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToAddGroupRequestCommand(req *request.AddGroup) *command.AddGroupRequest {
	return &command.AddGroupRequest{
		GroupID:   req.GroupID,
		Greeting:  req.Greeting,
		InvitedID: req.InvitedID,
		ApplyType: req.ApplyType,
	}
}

func ToListGroupRequestQuery(req *request.ListGroup) *query.ListGroupRequest {
	return &query.ListGroupRequest{
		FilterStatus: req.FilterStatus,
	}
}

func ToHandleGroupRequestCommand(req *request.HandleGroup) *command.HandleGroupRequest {
	return &command.HandleGroupRequest{
		HandleID: req.ID,
		Action:   req.Action,
	}
}

func ToListGroupRequestResult(row *model.GroupRequestRow) *response.ListGroupRequestResult {
	result := &response.ListGroupRequestResult{
		ID:             row.PublicID,
		GroupID:        row.GroupInfo.PublicID,
		GroupName:      row.GroupInfo.GroupName,
		GroupAvatar:    row.GroupInfo.Avatar,
		FromUserID:     row.ApplyUser.PublicID,
		FromUserName:   row.ApplyUser.Nickname,
		FromUserAvatar: row.ApplyUser.Avatar,
		Message:        row.Message,
		ApplyType:      row.GetApplyTypeText(),
		Status:         row.GetStatusText(),
		CreatedAt:      row.CreatedAt,
	}

	if row.InviterUser != nil {
		result.InviteUserID = row.InviterUser.PublicID
		result.InviteUserName = row.InviterUser.Nickname
		result.InviteUserAvatar = row.InviterUser.Avatar
	}

	if row.HandleByUser != nil {
		result.HandleByUserID = row.HandleByUser.PublicID
		result.HandleByUserName = row.HandleByUser.Nickname
		result.HandleByUserAvatar = row.HandleByUser.Avatar
	}

	return result
}
