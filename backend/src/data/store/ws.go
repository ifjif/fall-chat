package store

import (
	"zzc/fall-chat/src/data/cache"
	"zzc/fall-chat/src/data/repos"
)

type WSStore interface {
	WSRepos() repos.WSRepositories
	WSCache() cache.WSCache
}
