package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/utils"
)

func ToAddFriendCommand(req *request.AddFriend) *command.AddFriend {
	return &command.AddFriend{
		TargetID:  req.TargetID,
		ApplyType: req.ApplyType,
		Greeting:  req.Greeting,
	}
}

func ToListFriendRequestQuery(req *request.ListFriend) *query.ListFriendRequest {
	return &query.ListFriendRequest{
		FilterStatus: req.FilterStatus,
	}
}

func ToHandleFriendRequestCommand(req *request.HandleFriend) *command.HandleFriendRequest {
	return &command.HandleFriendRequest{
		HandleID: req.ID,
		Action:   req.Action,
	}
}

func ToListFriendRequestResult(m *model.SearchFriendRequestRow) *response.ListFriendRequestResult {
	return &response.ListFriendRequestResult{
		ID:         m.PublicID,
		FromUserID: utils.SF.Encode(m.FromUserID),
		ToUserID:   utils.SF.Encode(m.ToUserID),
		Message:    m.Message,
		Status:     m.GetStatusText(),
		ApplyType:  m.GetApplyTypeText(),
		Avatar:     m.Avatar,
		Nickname:   m.Nickname,
		CreatedAt:  m.CreatedAt,
	}
}
