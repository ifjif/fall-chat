package store

import (
	"context"
	"zzc/fall-chat/src/data/repos"
)

// 全dao store
// Store 是 service 唯一需要依赖的数据层接口。
// 它屏蔽了 gorm/redis 的细节，也方便单测时替换。
type Store interface {
	// 普通查询用
	Repos() repos.Repositories
	// 事务用
	WithinTX(ctx context.Context, fn func(repo repos.Repositories) error) error
}
