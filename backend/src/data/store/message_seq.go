package store

import (
	"context"
	"zzc/fall-chat/src/data/cache"
	"zzc/fall-chat/src/data/repos"
)

type MessageSeqStore interface {
	MessageSeqRepos() repos.MessageSeqRespositories
	WithinMessageSeqTX(ctx context.Context, fn func(repo repos.MessageSeqRespositories) error) error
	MessageSeqCache() cache.MessageSeqCache
}
