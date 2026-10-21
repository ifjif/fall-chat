package store

import (
	"context"
	"zzc/fall-chat/src/data/repos"
)

// group 业务store
type GroupStore interface {
	GroupRepos() repos.GroupRepositories
	WithinGroupTX(ctx context.Context, fn func(repo repos.GroupRepositories) error) error
}
