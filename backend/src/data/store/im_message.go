package store

import "zzc/fall-chat/src/data/repos"

type IMMessageStore interface {
	IMMessageRepos() repos.IMMessageRepositories
}
