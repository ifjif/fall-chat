package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

// group专属
type GroupRepositories struct {
	Group        dao.GroupDao
	GroupMember  dao.GroupMemberDao
	Conversation dao.ConversationDao
	ConvMemb     dao.ConvMembDao
}

func NewGroupRepos(e executor.Executor) GroupRepositories {
	return GroupRepositories{
		Group:        dao.NewGroupDao(e),
		GroupMember:  dao.NewGroupMemberDao(e),
		Conversation: dao.NewConversationDao(e),
		ConvMemb:     dao.NewConvMembDao(e),
	}
}
