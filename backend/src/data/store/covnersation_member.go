package store

import (
	"context"
	"zzc/fall-chat/src/data/cache"
	"zzc/fall-chat/src/data/repos"
)

type ConvMembStore interface {
	ConvMembRepos() repos.ConvMembRepositories
	WithinConvMembTX(ctx context.Context, fn func(repo repos.ConvMembRepositories) error) error
	ConvMembCache() cache.ConvMembCache
}
