package service

import (
	"context"
	"encoding/json"
	"time"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/job"
	"zzc/fall-chat/src/mess"
	messagequeue "zzc/fall-chat/src/message_queue"
	"zzc/fall-chat/src/query"
	taskmanager "zzc/fall-chat/src/task_manager"

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
	messageGroup := messagequeue.GetDispatcher().GetOrCreateGroup("update-conv-member", "update-conv-member-group")
	messageGroup.RegisterHandler(cs.UpdateConversationHandler)
	messageGroup.StartInstance()
}

func (cs *ConversationService) List(id int64, name string) ([]*response.IMConversationResponse, error) {
	var result []*response.IMConversationResponse
	rows, err := cs.dao.ListByUserID(id)
	if err != nil {
		return nil, cs.Wrap("找用户会话列表", err, "id", id, "username", name)
	}

	for i := range rows {
		row := rows[i]
		v := converter.ToIMConversationResponse(&row)
		v.PeerID = cs.sf.Encode(row.PeerID)
		result = append(result, v)
	}
	return result, err
}

func (cs *ConversationService) Single(qu *query.ConvSingleQuery) (*response.IMConversationResponse, error) {
	convID, err := cs.sf.Decode(qu.ConvID)
	if err != nil {
		return nil, cs.Wrap("会话id解码", err, "id", qu.ID)
	}

	row, err := cs.dao.SingleByConvID(convID, qu.ID)
	if err != nil {
		return nil, cs.Wrap("找用户的单条会话", err, "id", qu.ID, "conv_id", convID)
	}

	// 没有会话
	if row == nil {
		return nil, nil
	}

	result := converter.ToIMConversationResponse(row)
	result.PeerID = cs.sf.Encode(row.PeerID)
	return result, nil
}

func (cs *ConversationService) UpdateConversationHandler(ctx context.Context, msg *kafka.Message, commitFunc func()) error {
	var updateMess mess.UpdateConvMess
	err := json.Unmarshal(msg.Value, &updateMess)
	if err != nil {
		return cs.Wrap("json unmarshal kafka消息", err)
	}

	task := job.NewUpdateConversaionJob(cs.db)
	task.ConvID = updateMess.ConvID
	task.PeerID = updateMess.PeerID
	task.Seq = updateMess.Seq
	task.Duration = 500 * time.Millisecond
	task.MsgID = updateMess.MsgID
	task.MsgTime = updateMess.MsgTime
	task.MsgPreview = updateMess.MsgPreview
	task.MsgSenderID = updateMess.MsgSenderID
	task.CompleteFunc = func() {
		commitFunc()
	}

	taskmanager.TaskManager.Submit(task)

	return nil
}
