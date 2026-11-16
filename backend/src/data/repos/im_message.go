package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type IMMessageRepositories struct {
	IMMessage dao.IMMessageDao
}

func NewIMMessageRepos(db executor.Executor) IMMessageRepositories {
	return IMMessageRepositories{
		IMMessage: dao.NewIMMessageDao(db),
	}
}
