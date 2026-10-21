package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

type GroupRequestRepositories struct {
	GroupRequest dao.GroupRequestDao
	Group        dao.GroupDao
	GroupMemeber dao.GroupMemberDao
	ConvMemb     dao.ConvMembDao
}

func NewGroupRequestRepos(db executor.Executor) GroupRequestRepositories {
	return GroupRequestRepositories{
		GroupRequest: dao.NewGroupRequestDao(db),
		Group:        dao.NewGroupDao(db),
		GroupMemeber: dao.NewGroupMemberDao(db),
		ConvMemb:     dao.NewConvMembDao(db),
	}
}
