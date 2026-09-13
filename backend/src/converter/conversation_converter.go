package converter

import (
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToIMConversationResponse(row *model.IMConversationRow) *response.IMConversationResponse {
	return &response.IMConversationResponse{
		ID:              row.PublicID,
		Type:            row.GetConvTypeText(),
		Name:            row.PeerNickname,
		Avatar:          row.PeerAvatar,
		LastMsgID:       row.LastMsgID,
		LastMsgPreview:  row.LastMsgPreview,
		LastMsgSenderID: row.LastMsgSenderID,
		LastMsgTime:     row.LastMsgTime,
		LastMsgSeq:      row.LastMsgSeq,
		LastReadSeq:     row.LastReadSeq,
		UnreadCount:     row.LastMsgSeq - row.LastReadSeq,
	}
}

func ToConvSingleQuery(req *request.ConvSingleRequest) *query.ConvSingleQuery {
	return &query.ConvSingleQuery{
		ConvID: req.ConvID,
	}
}
