package data

import (
	"context"
	"sync"
	"zzc/fall-chat/src/data/cache"
	"zzc/fall-chat/src/data/repos"
	"zzc/fall-chat/src/database"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	bd   *BaseDao
	once sync.Once
)

type BaseDao struct {
	db  *gorm.DB
	rdb *redis.Client
}

func GetBaseDao() *BaseDao {
	once.Do(func() {
		bd = &BaseDao{
			db:  database.DB,
			rdb: database.Redis(),
		}
	})

	return bd
}

// ======================================整仓库
// 返回非事务作用域的一套仓库
func (d *BaseDao) Repos() repos.Repositories {
	return repos.NewRepos(d.db)
}

// 事务执行器，service只拿到repositories
// 事务边界由调用它的service决定
func (d *BaseDao) WithinTX(ctx context.Context, fn func(repos.Repositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 事务作用域内的仓库，所有DAO 都绑定在tx上
		return fn(repos.NewRepos(tx))
	})
}

// ======================================group仓库
func (d *BaseDao) GroupRepos() repos.GroupRepositories {
	return repos.NewGroupRepos(d.db)
}

func (d *BaseDao) WithinGroupTX(ctx context.Context, fn func(repos.GroupRepositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repos.NewGroupRepos(tx))
	})
}

// ======================================conversation仓库
func (d *BaseDao) ConvRepos() repos.ConvRepositories {
	return repos.NewConvRepos(d.db)
}

func (d *BaseDao) WithinConvTX(ctx context.Context, fn func(repo repos.ConvRepositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repos.NewConvRepos(tx))
	})
}

func (d *BaseDao) ConvCache() cache.ConversationCache {
	return cache.NewConversationCache(d.rdb)
}

// ======================================conversation member仓库
func (d *BaseDao) ConvMembRepos() repos.ConvMembRepositories {
	return repos.NewConvMembRepos(d.db)
}

func (d *BaseDao) WithinConvMembTX(ctx context.Context, fn func(repo repos.ConvMembRepositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repos.NewConvMembRepos(tx))
	})
}

func (d *BaseDao) ConvMembCache() cache.ConvMembCache {
	return cache.NewConvMembCache(d.rdb)
}

// ======================================email仓库

func (d *BaseDao) EmailRepos() repos.EmailRepositories {
	return repos.NewEmailRepos(d.db)
}

func (d *BaseDao) WithinEmailTX(ctx context.Context, fn func(repo repos.EmailRepositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repos.NewEmailRepos(tx))
	})
}

// ======================================file仓库
func (d *BaseDao) FileRepos() repos.FileRepositories {
	return repos.NewFileRepos(d.db)
}

func (d *BaseDao) WithinFileTX(ctx context.Context, fn func(repo repos.FileRepositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repos.NewFileRepos(tx))
	})
}

// ======================================message_seq仓库

func (d *BaseDao) MessageSeqRepos() repos.MessageSeqRespositories {
	return repos.NewMessageSeqRepos(d.db)
}

func (d *BaseDao) WithinMessageSeqTX(ctx context.Context, fn func(repo repos.MessageSeqRespositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repos.NewMessageSeqRepos(tx))
	})
}

func (d *BaseDao) MessageSeqCache() cache.MessageSeqCache {
	return cache.NewMessageSeqCache(d.rdb)
}

// ======================================user仓库
func (d *BaseDao) UserRepos() repos.UserRepositories {
	return repos.NewUserRepos(d.db)
}

func (d *BaseDao) WithinUserTX(ctx context.Context, fn func(repo repos.UserRepositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repos.NewUserRepos(tx))
	})
}

// ======================================search仓库
func (d *BaseDao) SearchRepos() repos.SearchRepositories {
	return repos.NewSearchRepos(d.db)
}

// ======================================friendship仓库
func (d *BaseDao) FriendshipRepos() repos.FriendshipRepositories {
	return repos.NewFriendshipRepos(d.db)
}

// ======================================friendship_request仓库
func (d *BaseDao) FriendshipRequestRepos() *repos.FriendshipRequestRepositories {
	return repos.NewFriendshipRequestRepositories(d.db)
}

func (d *BaseDao) WithinFriendshipRequestTX(ctx context.Context, fn func(repo repos.FriendshipRequestRepositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(*repos.NewFriendshipRequestRepositories(tx))
	})
}

// ======================================group_request仓库
func (d *BaseDao) GroupRequestRepos() repos.GroupRequestRepositories {
	return repos.NewGroupRequestRepos(d.db)
}

func (d *BaseDao) WithinGroupRequestTX(ctx context.Context, fn func(repo repos.GroupRequestRepositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repos.NewGroupRequestRepos(tx))
	})
}

// ======================================ws仓库
func (d *BaseDao) WSRepos() repos.WSRepositories {
	return repos.NewWSRepos(d.db)
}

func (d *BaseDao) WSCache() cache.WSCache {
	return cache.NewWSCache(d.rdb)
}

// ======================================group_member仓库
func (d *BaseDao) GroupMemberRepos() repos.GroupMemberRepositories {
	return repos.NewGroupMemberRepo(d.db)
}

func (d *BaseDao) GroupMemberCache() cache.GroupMemberCache {
	return cache.NewGroupMemberCache(d.rdb)
}

// ======================================update_conv_job仓库
func (d *BaseDao) WithinUpdateConvJobTX(ctx context.Context, fn func(repo repos.UpdateConvJobRepositories) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(repos.NewUpdateConvJobRepositories(tx))
	})
}
