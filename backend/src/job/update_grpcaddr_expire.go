package job

import (
	"context"
	"time"
	"zzc/fall-chat/src/data/cache"
)

type UpdateGrpcaddrExpire struct {
	Common
	Duration time.Duration
	ID       string
	cache    cache.WSCache
}

func NewUpdateGrpcaddrExpire(cache cache.WSCache) *UpdateGrpcaddrExpire {
	return &UpdateGrpcaddrExpire{
		cache: cache,
	}
}

func (u *UpdateGrpcaddrExpire) Key() string {
	return u.ID
}

func (u *UpdateGrpcaddrExpire) Delay() time.Duration {
	return u.Duration
}

func (u *UpdateGrpcaddrExpire) Execute(ctx context.Context) error {
	u.cache.AddGrpcAddr(ctx, u.ID)
	return nil
}
