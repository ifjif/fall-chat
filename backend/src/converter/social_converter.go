package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/utils"
)

func ToSocialAddFriendCommand(req *request.SocialAddFriendRequest) *command.SocialAddFriendCommand {
	return &command.SocialAddFriendCommand{
		TargetID:  req.TargetID,
		ApplyType: req.ApplyType,
		Greeting:  req.Greeting,
	}
}

func ToSocialSearchFriendRequestQuery(req *request.SocialSearchFriendRequestRequest) *query.SocialSearchFriendRequestQuery {
	return &query.SocialSearchFriendRequestQuery{
		FilterStatus: req.FilterStatus,
	}
}

func ToSocialSearchFriendRequestResult(m *model.SearchFriendRequestRow) *response.SocialSearchFriendRequestResult {
	return &response.SocialSearchFriendRequestResult{
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

func ToSocialHandleFriendRequestCommand(req *request.SocialHandleFriendRequestRequest) *command.SocialHandleFriendRequestCommand {
	return &command.SocialHandleFriendRequestCommand{
		HandleID: req.ID,
		Action:   req.Action,
	}
}

func ToSocialAddGroupRequestCommand(req *request.SocialAddGroupRequestRequest) *command.SocialAddGroupRequestCommand {
	return &command.SocialAddGroupRequestCommand{
		GroupID:   req.GroupID,
		Greeting:  req.Greeting,
		InvitedID: req.InvitedID,
		ApplyType: req.ApplyType,
	}
}

func ToSocialSearchGroupRequestQuery(req *request.SocialSearchGroupRequestRequest) *query.SocialSearchGroupRequestQuery {
	return &query.SocialSearchGroupRequestQuery{
		FilterStatus: req.FilterStatus,
	}
}

func ToSocialSearchGroupRequestResult(row *model.GroupRequestRow) *response.SocialSearchGroupRequestResult {
	result := &response.SocialSearchGroupRequestResult{
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

func ToSocialHandleGroupRequestCommand(req *request.SocialHandleGroupRequestRequest) *command.SocialHandleGroupRequestCommand {
	return &command.SocialHandleGroupRequestCommand{
		HandleID: req.ID,
		Action:   req.Action,
	}
}
