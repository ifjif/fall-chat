package service

import (
	"context"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
)

type Friendship struct {
	BaseService
	store store.FriendshipStore
}

func NewFriendship(name string) *Friendship {
	bs := BaseService{name: name}
	return &Friendship{
		BaseService: bs,
	}
}

func (fs *Friendship) Inject() {
	fs.BaseService.Inject()
	fs.store = data.GetBaseDao()
}

func (fs *Friendship) FriendshipList(ctx context.Context, qu *query.FriendshipList) ([]model.FriendshipRow, error) {
	rows, err := fs.store.
		FriendshipRepos().
		Friendship.
		FriendshipList(ctx, qu.ID, 10)
	if err != nil {
		return nil, fs.Wrap("查找好友列表", err, "username", qu.Username)
	}
	return rows, nil
}
