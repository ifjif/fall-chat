package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type EmailRepositories struct {
	Email dao.EmailDao
}

func NewEmailRepos(db executor.Executor) EmailRepositories {
	return EmailRepositories{
		Email: dao.NewEmailDao(db),
	}
}
