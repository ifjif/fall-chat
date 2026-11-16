package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/filter"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToAddGroupRequestCommand(req *request.AddGroup) (*command.AddGroupRequest, error) {
	groupID, err := model.DecodeID(req.GroupID)
	if err != nil {
		return nil, err
	}
	var invitedID int64
	if req.InvitedID != "" {
		invitedID, err = model.DecodeID(req.InvitedID)
		if err != nil {
			return nil, err
		}
	}
	cmd := &command.AddGroupRequest{
		GroupID:   groupID,
		Greeting:  req.Greeting,
		ApplyType: req.ApplyType,
	}
	if req.InvitedID != "" {
		cmd.InvitedID = &invitedID
	}

	return cmd, nil
}

func ToListGroupRequestQuery(req *request.ListGroup) *query.ListGroupRequest {
	qu := &query.ListGroupRequest{
		Next:    req.Next,
		SortKey: req.SortKey,
	}

	qu.Filter = &filter.GroupRequestFilter{
		Status: req.Status,
	}

	return qu
}

func ToHandleGroupRequestCommand(req *request.HandleGroup) (*command.HandleGroupRequest, error) {
	handleID, err := model.DecodeID(req.ID)
	if err != nil {
		return nil, err
	}
	return &command.HandleGroupRequest{
		HandleID: handleID,
		Action:   req.Action,
	}, nil
}

func ToGroupRequestResult(row *model.GroupRequestRow, userID uint64) *response.GroupRequestResult {
	result := &response.GroupRequestResult{
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

	result.Sender = row.ApplyUser.ID == userID

	return result
}
