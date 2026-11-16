package converter

import (
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToIMConversationResponse(row *model.IMConversationRow) *response.IMConversation {
	return &response.IMConversation{
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

func ToConvSingleQuery(req *request.ConvSingle) (*query.ConvSingle, error) {
	convID, err := model.DecodeID(req.ConvID)
	if err != nil {
		return nil, err
	}
	return &query.ConvSingle{
		ConvID: convID,
	}, nil
}
