package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type ConvRepositories struct {
	Conversation dao.ConversationDao
	ConvMemb     dao.ConvMembDao
	User         dao.UserDao
}

func NewConvRepos(db executor.Executor) ConvRepositories {
	return ConvRepositories{
		Conversation: dao.NewConversationDao(db),
		ConvMemb:     dao.NewConvMembDao(db),
		User:         dao.NewUserDao(db),
	}
}
