package service

import (
	"context"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/dto/response"
	messagequeue "zzc/fall-chat/src/message_queue"

	"github.com/segmentio/kafka-go"
)

type ConversationService struct {
	BaseService
	dao *data.ConversationDao
}

func NewConversationService(name string) *ConversationService {
	bs := BaseService{name: name}
	return &ConversationService{
		BaseService: bs,
	}
}

func (cs *ConversationService) Inject() {
	cs.BaseService.Inject()
	cs.dao = data.NewConversationDao(cs.db)
	messageGroup := messagequeue.GetDispatcher().GetOrCreateGroup("update-conv&member", "update-conv&member-group")
	messageGroup.RegisterHandler(cs.UpdateConversationHandler)
}

func (cs *ConversationService) List(id int64, name string) ([]*response.IMConversationListResponse, error) {
	var result []*response.IMConversationListResponse
	rows, err := cs.dao.ListByUserID(id)
	if err != nil {
		cs.LogWarn("List", err, "username", name, "reason", "找用户会话列表失败")
		return nil, err
	}

	for i := range rows {
		row := rows[i]
		v := converter.ToIMConversationListResponse(&row)
		v.PeerID = cs.sf.Encode(row.PeerID)
		result = append(result, v)
	}
	return result, err
}

func (cs *ConversationService) UpdateConversationHandler(ctx context.Context, msg *kafka.Message, commitFunc func()) error {
	return nil
}
