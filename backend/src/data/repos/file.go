package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type FileRepositories struct {
	File dao.FileDao
}

func NewFileRepos(db executor.Executor) FileRepositories {
	return FileRepositories{
		File: dao.NewFileDao(db),
	}
}
