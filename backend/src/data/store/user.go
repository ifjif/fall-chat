package store

import (
	"context"
	"zzc/fall-chat/src/data/repos"
)

type UserStore interface {
	UserRepos() repos.UserRepositories
	WithinUserTX(ctx context.Context, fn func(repo repos.UserRepositories) error) error
}
