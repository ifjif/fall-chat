package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type GroupMemberRepositories struct {
	GroupMember dao.GroupMemberDao
}

func NewGroupMemberRepo(db executor.Executor) GroupMemberRepositories {
	return GroupMemberRepositories{
		GroupMember: dao.NewGroupMemberDao(db),
	}
}
