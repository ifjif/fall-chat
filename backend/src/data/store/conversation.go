package store

import (
	"context"
	"zzc/fall-chat/src/data/cache"
	"zzc/fall-chat/src/data/repos"
)

type ConversationStore interface {
	ConvRepos() repos.ConvRepositories
	WithinConvTX(ctx context.Context, fn func(repo repos.ConvRepositories) error) error
	ConvCache() cache.ConversationCache
}
