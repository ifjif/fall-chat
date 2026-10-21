package store

import (
	"zzc/fall-chat/src/data/cache"
	"zzc/fall-chat/src/data/repos"
)

type GroupMemberStore interface {
	GroupMemberRepos() repos.GroupMemberRepositories
	GroupMemberCache() cache.GroupMemberCache
}
