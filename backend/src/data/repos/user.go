package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type UserRepositories struct {
	User dao.UserDao
}

func NewUserRepos(db executor.Executor) UserRepositories {
	return UserRepositories{
		User: dao.NewUserDao(db),
	}
}
