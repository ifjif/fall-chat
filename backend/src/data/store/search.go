package store

import "zzc/fall-chat/src/data/repos"

type SearchStore interface {
	SearchRepos() repos.SearchRepositories
}
