package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type FriendshipRepositories struct {
	Friendship dao.FriendshipDao
}

func NewFriendshipRepos(db executor.Executor) FriendshipRepositories {
	return FriendshipRepositories{
		Friendship: dao.NewFriendshipDao(db),
	}
}
