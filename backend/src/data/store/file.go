package store

import (
	"context"
	"zzc/fall-chat/src/data/repos"
)

type FileStore interface {
	FileRepos() repos.FileRepositories
	WithinFileTX(ctx context.Context, fn func(repo repos.FileRepositories) error) error
}
