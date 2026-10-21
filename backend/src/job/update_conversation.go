package job

import (
	"context"
	"fmt"
	"time"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/repos"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/model"
	taskmanager "zzc/fall-chat/src/task_manager"
)

type UpdateConversation struct {
	Common
	ConvID       int64 // 会话id
	PeerID       int64 // 对方id 用户或群
	Seq          int64 // 消息序列
	ConvType     int   // 会话类型 0:单聊  1:群
	MsgID        int64
	MsgPreview   string
	MsgSenderID  int64
	MsgTime      time.Time
	Duration     time.Duration
	CompleteFunc func()
	store        store.UpdateConvJobStore
}

func NewUpdateConversaion() *UpdateConversation {
	return &UpdateConversation{
		store: data.GetBaseDao(),
	}
}

func (uc *UpdateConversation) Key() string {
	return fmt.Sprintf("im-conv-key-%d", uc.ConvID)
}

func (uc *UpdateConversation) Delay() time.Duration {
	return uc.Duration
}

func (uc *UpdateConversation) ShouldReplace(old taskmanager.Job) bool {
	if ouc, ok := old.(*UpdateConversation); ok {
		return uc.Seq > ouc.Seq
	}

	return true
}

func (uc *UpdateConversation) Execute(ctx context.Context) error {
	updateConv := &model.IMConversation{}
	updateConv.ID = uint64(uc.ConvID)
	updateConv.LastMsgID = uc.MsgID
	updateConv.LastMsgTime = uc.MsgTime
	updateConv.LastMsgPreview = uc.MsgPreview
	updateConv.LastMsgSenderID = uc.MsgSenderID
	updateConv.LastMsgSeq = uc.Seq

	updateMember := &model.IMConversationMember{}
	updateMember.ConvID = uc.ConvID
	updateMember.UserID = uc.MsgSenderID
	updateMember.LastReadMsgID = uc.MsgID
	updateMember.LastReadSeq = uc.Seq
	reason := ""
	err := uc.store.WithinUpdateConvJobTX(ctx, func(repo repos.UpdateConvJobRepositories) error {
		reason = "更新会话最后一条消息失败"
		// 更新会话最后一条消息
		err := repo.Convsation.UpdateLastMessage(ctx, updateConv)
		if err != nil {
			return err
		}

		reason = "更新会话发送者消息信息失败"
		// 更新成员消息seq
		err = repo.ConvMemb.UpdateLastReadSeq(ctx, updateMember)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		uc.LogWarn("UpdateConversationJob", err,
			"conv_id", uc.ConvID, "msg_sender_id", uc.MsgSenderID, "reason", reason)

		return err
	}
	return nil
}

func (uc *UpdateConversation) OnComplete() {
	if uc.CompleteFunc != nil {
		uc.CompleteFunc()
	}
}
