package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type UpdateConvJobRepositories struct {
	Convsation dao.ConversationDao
	ConvMemb   dao.ConvMembDao
}

func NewUpdateConvJobRepositories(db executor.Executor) UpdateConvJobRepositories {
	return UpdateConvJobRepositories{
		Convsation: dao.NewConversationDao(db),
		ConvMemb:   dao.NewConvMembDao(db),
	}
}
