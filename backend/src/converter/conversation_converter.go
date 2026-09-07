package converter

import (
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
)

func ToIMConversationListResponse(row *model.IMConversationRow) *response.IMConversationListResponse {
	return &response.IMConversationListResponse{
		ID:              row.PublicID,
		Type:            row.GetConvTypeText(),
		Name:            row.PeerNickname,
		Avatar:          row.PeerAvatar,
		LastMsgID:       row.LastMsgID,
		LastMsgPreview:  row.LastMsgPreview,
		LastMsgSenderID: row.LastMsgSenderID,
		LastMsgTime:     row.LastMsgTime,
	}
}
