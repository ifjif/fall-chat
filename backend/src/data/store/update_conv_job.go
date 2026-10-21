package store

import (
	"context"
	"zzc/fall-chat/src/data/repos"
)

type UpdateConvJobStore interface {
	WithinUpdateConvJobTX(ctx context.Context, fn func(repo repos.UpdateConvJobRepositories) error) error
}
