package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type ConvMembRepositories struct {
	ConvMemb dao.ConvMembDao
}

func NewConvMembRepos(db executor.Executor) ConvMembRepositories {
	return ConvMembRepositories{
		ConvMemb: dao.NewConvMembDao(db),
	}
}
