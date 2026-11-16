package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type EmojiRepositories struct {
	EmojiPack dao.EmojiPackDao
	EmojiItem dao.EmojiItemDao
	EmojiOwn  dao.EmojiOwnDao
}

func NewEmojiRepos(db executor.Executor) EmojiRepositories {
	return EmojiRepositories{
		EmojiPack: dao.NewEmojiPackDao(db),
		EmojiItem: dao.NewEmojiItemDao(db),
		EmojiOwn:  dao.NewEmojiOwnDao(db),
	}
}
