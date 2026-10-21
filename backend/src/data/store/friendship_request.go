package store

import (
	"context"
	"zzc/fall-chat/src/data/repos"
)

type FriendshipRequestStore interface {
	FriendshipRequestRepos() *repos.FriendshipRequestRepositories
	WithinFriendshipRequestTX(ctx context.Context, fn func(repo repos.FriendshipRequestRepositories) error) error
}
