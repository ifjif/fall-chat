package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type FriendshipRequestRepositories struct {
	FriendshipRequest dao.FriendshipRequestDao
	Friendship        dao.FriendshipDao
}

func NewFriendshipRequestRepositories(db executor.Executor) *FriendshipRequestRepositories {
	return &FriendshipRequestRepositories{
		FriendshipRequest: dao.NewFriendshipRequestDao(db),
		Friendship:        dao.NewFriendshipDao(db),
	}
}
