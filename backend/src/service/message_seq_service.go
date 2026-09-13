package service

import "zzc/fall-chat/src/data"

type MessageSeqService struct {
	BaseService
	dao *data.MessageSeqDao
}

func NewMessageSeqService(name string) *MessageSeqService {
	bs := BaseService{name: name}

	return &MessageSeqService{
		BaseService: bs,
	}
}

func (mss *MessageSeqService) Inject() {
	mss.BaseService.Inject()
	dao := data.NewMessageSeqDao(mss.db)
	mss.dao = dao
}

func (mss *MessageSeqService) NextSeq(convID int64) (int64, error) {
	v, err := mss.dao.GenerateNextSeq(convID)
	if err != nil {
		return 0, mss.LogWarn("NextSeq", err, "convID", convID,
			"reason", "得到当前会话消息的下一个seq失败")
	}

	return v, nil
}
