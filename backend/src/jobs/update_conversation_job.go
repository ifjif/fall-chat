package jobs

import (
	"fmt"
	"time"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/model"
	taskmanager "zzc/fall-chat/src/task_manager"

	"gorm.io/gorm"
)

type UpdateConversationJob struct {
	CommonJob
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
	dao          *data.UpdateConversationJobDao
}

func NewUpdateConversaionJob(db *gorm.DB) *UpdateConversationJob {
	return &UpdateConversationJob{
		dao: data.GetUpdateConversationJobDao(db),
	}
}

func (uc *UpdateConversationJob) Key() string {
	return fmt.Sprintf("im-conv-key-%d", uc.ConvID)
}

func (uc *UpdateConversationJob) Delay() time.Duration {
	return uc.Duration
}

func (uc *UpdateConversationJob) ShouldReplace(old taskmanager.Job) bool {
	if ouc, ok := old.(*UpdateConversationJob); ok {
		return uc.Seq > ouc.Seq
	}

	return true
}

func (uc *UpdateConversationJob) Execute() error {
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
	err := uc.dao.DB().Transaction(func(tx *gorm.DB) error {
		reason = "更新会话最后一条消息失败"
		// 更新会话最后一条消息
		err := uc.dao.UpdateConvLastMessWithTx(tx, updateConv)
		if err != nil {
			return err
		}

		reason = "更新会话发送者消息信息失败"
		// 更新成员消息seq
		err = uc.dao.UpdateConvMemberMessSeqWithTx(tx, updateMember)
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

func (uc *UpdateConversationJob) OnComplete() {
	if uc.CompleteFunc != nil {
		uc.CompleteFunc()
	}
}
