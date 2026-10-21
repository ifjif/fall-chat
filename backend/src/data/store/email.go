package store

import (
	"context"
	"zzc/fall-chat/src/data/repos"
)

type EmailStore interface {
	EmailRepos() repos.EmailRepositories
	WithinEmailTX(ctx context.Context, fn func(repo repos.EmailRepositories) error) error
}
