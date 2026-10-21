package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type SearchRepositories struct {
	User  dao.UserDao
	Group dao.GroupDao
}

func NewSearchRepos(db executor.Executor) SearchRepositories {
	return SearchRepositories{
		User:  dao.NewUserDao(db),
		Group: dao.NewGroupDao(db),
	}
}
