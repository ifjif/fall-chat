package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type MessageSeqRespositories struct {
	MessageSeq dao.MessageSeqDao
}

func NewMessageSeqRepos(db executor.Executor) MessageSeqRespositories {
	return MessageSeqRespositories{
		MessageSeq: dao.NewMessageSeqDao(db),
	}
}
