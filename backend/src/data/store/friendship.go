package store

import "zzc/fall-chat/src/data/repos"

type FriendshipStore interface {
	FriendshipRepos() repos.FriendshipRepositories
}
