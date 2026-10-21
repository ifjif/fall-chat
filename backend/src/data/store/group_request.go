package store

import (
	"context"
	"zzc/fall-chat/src/data/repos"
)

type GroupRequestStore interface {
	GroupRequestRepos() repos.GroupRequestRepositories
	WithinGroupRequestTX(ctx context.Context, fn func(repo repos.GroupRequestRepositories) error) error
}
