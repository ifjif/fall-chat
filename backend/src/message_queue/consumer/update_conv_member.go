package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"zzc/fall-chat/src/job"
	"zzc/fall-chat/src/mess"
	messagequeue "zzc/fall-chat/src/message_queue"
	taskmanager "zzc/fall-chat/src/task_manager"

	"github.com/segmentio/kafka-go"
)

var (
	updateConvMemberOnce sync.Once
	updateConvMember     *UpdateConvMember
)

type UpdateConvMember struct {
	writer *kafka.Writer
}

func (h *UpdateConvMember) init() {
	messageGroup := messagequeue.GetDispatcher().GetOrCreateGroup("update-conv-member", "update-conv-member-group")
	messageGroup.RegisterHandler(h.Handle)
	messageGroup.StartInstance()
}

func GetUpdateConvMember() *UpdateConvMember {
	if updateConvMember != nil {
		return updateConvMember
	}

	updateConvMemberOnce.Do(func() {
		updateConvMember = &UpdateConvMember{}
		updateConvMember.init()
	})

	return updateConvMember
}

func (h *UpdateConvMember) Handle(ctx context.Context, msg *kafka.Message, commitFunc func()) error {
	var updateMess mess.UpdateConv
	err := json.Unmarshal(msg.Value, &updateMess)
	if err != nil {
		commitFunc()
		// 脏消息, 丢弃
		return fmt.Errorf("json unmarshal kafka消息", err)
	}

	jo := job.NewUpdateConversaion()
	jo.ConvID = updateMess.ConvID
	jo.PeerID = updateMess.PeerID
	jo.Seq = updateMess.Seq
	jo.Duration = 500 * time.Millisecond
	jo.MsgID = updateMess.MsgID
	jo.MsgTime = updateMess.MsgTime
	jo.MsgPreview = updateMess.MsgPreview
	jo.MsgSenderID = updateMess.MsgSenderID
	jo.CompleteFunc = func() {
		commitFunc()
	}

	taskmanager.TaskManager.Submit(jo)

	return nil
}
