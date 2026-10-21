package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type WSRepositories struct {
	IMMessage dao.IMMessageDao
}

func NewWSRepos(db executor.Executor) WSRepositories {
	return WSRepositories{
		IMMessage: dao.NewIMMessageDao(db),
	}
}
