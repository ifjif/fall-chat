package converter

import (
	"zzc/fall-chat/src/dto/request"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

func ToPatchListIMMessageQuery(req *request.PatchListIMMessage) (*query.PatchListIMMessage, error) {
	convID, err := model.DecodeID(req.ConvID)
	if err != nil {
		return nil, err
	}
	return &query.PatchListIMMessage{
		ConvID: convID,
		MinSeq: *req.MinSeq,
		MaxSeq: *req.MaxSeq,
	}, nil
}

func ToListIMMessageQuery(req *request.ListIMMessage) (*query.ListIMMessage, error) {
	convID, err := model.DecodeID(req.ConvID)
	if err != nil {
		return nil, err
	}
	return &query.ListIMMessage{
		ConvID:  convID,
		Next:    req.Next,
		SortKey: req.SortKey,
	}, nil
}

func ToImMessageResult(m *model.IMMessage) *response.IMMessageResult {
	return &response.IMMessageResult{
		ID:       m.PublicID,
		ConvID:   m.EncodeConvID(),
		Content:  m.Content,
		SenderID: m.EncodeSenderID(),
		SendTime: m.SendTime,
		MsgType:  m.MsgType,
		Status:   m.GetStatusText(),
		Seq:      m.Seq,
	}
}
