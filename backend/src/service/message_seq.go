package service

import (
	"context"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/keys/ctxkey"
)

type MessageSeq struct {
	BaseService
	store store.MessageSeqStore
}

func NewMessageSeq(name string) *MessageSeq {
	bs := BaseService{name: name}

	return &MessageSeq{
		BaseService: bs,
	}
}

func (mss *MessageSeq) Inject() {
	mss.BaseService.Inject()
	mss.store = data.GetBaseDao()
}

func (mss *MessageSeq) NextSeq(ctx context.Context, convID int64) (int64, error) {
	logger := ctxkey.GetLogger(ctx)
	// 从缓存中找
	seq, ok := mss.store.
		MessageSeqCache().
		IncrConvSeq(ctx, convID)
	if ok {
		return seq, nil
	}

	// 没有从数据库
	seq, err := mss.store.
		MessageSeqRepos().
		MessageSeq.
		GenerateNextSeq(ctx, convID)
	if err != nil {
		return 0, mss.Wrap("得到当前会话消息的下一个seq", err, "convID", convID)
	}

	// 设置缓存
	err = mss.store.
		MessageSeqCache().
		SetConvSeq(ctx, convID, seq)
	if err != nil {
		logger.Error("[redis]",
			"reason", "设置会话消息seq失败",
			"conv_id", convID,
			"seq", seq,
			"error", err)
	}

	return seq, nil
}
