package store

import (
	"context"
	"zzc/fall-chat/src/data/cache"
	"zzc/fall-chat/src/data/repos"
)

type EmojiStore interface {
	EmojiRepos() repos.EmojiRepositories
	WithinEmojiTX(ctx context.Context, fn func(repo repos.EmojiRepositories) error) error
	EmojiCache() cache.EmojiCache
}
