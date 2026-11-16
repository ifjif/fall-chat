package converter

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/filter"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToAddFriendCommand(req *request.AddFriend) (*command.AddFriend, error) {
	targetID, err := model.DecodeID(req.TargetID)
	if err != nil {
		return nil, err
	}
	return &command.AddFriend{
		TargetID:  targetID,
		ApplyType: req.ApplyType,
		Greeting:  req.Greeting,
	}, nil
}

func ToListFriendRequestQuery(req *request.ListFriend) *query.ListFriendRequest {
	qu := &query.ListFriendRequest{
		Next:    req.Next,
		SortKey: req.SortKey,
	}

	qu.Filter = &filter.FriendshipRequestFilter{
		Status: req.Status,
	}

	return qu
}

func ToHandleFriendRequestCommand(req *request.HandleFriend) (*command.HandleFriendRequest, error) {
	handleID, err := model.DecodeID(req.ID)
	if err != nil {
		return nil, err
	}
	return &command.HandleFriendRequest{
		HandleID: handleID,
		Action:   req.Action,
	}, nil
}

func ToFriendRequestResult(m *model.FriendshipRequestRow, userID int64) *response.FriendRequestResult {
	r := &response.FriendRequestResult{
		ID:         m.PublicID,
		FromUserID: model.EncodeID(m.FromUserID),
		ToUserID:   model.EncodeID(m.ToUserID),
		Message:    m.Message,
		Status:     m.GetStatusText(),
		ApplyType:  m.GetApplyTypeText(),
		Avatar:     m.Avatar,
		Nickname:   m.Nickname,
		CreatedAt:  m.CreatedAt,
	}

	r.Sender = m.FromUserID == userID

	return r
}
