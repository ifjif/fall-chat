package repos

import (
	"zzc/fall-chat/src/data/dao"
	"zzc/fall-chat/src/data/executor"
)

// 把业务用到的dao聚到一起
// service层永远是“绑定了某个executor(普通连接or事务)”的一套仓库
type Repositories struct {
	Group        dao.GroupDao
	GroupMember  dao.GroupMemberDao
	Conversation dao.ConversationDao
	ConvMemb     dao.ConvMembDao
}

// 唯一构造入口
// 传d.db -> 非事务仓库；
// 传tx -> 事务内仓库;
func NewRepos(e executor.Executor) Repositories {
	return Repositories{
		Group:        dao.NewGroupDao(e),
		GroupMember:  dao.NewGroupMemberDao(e),
		Conversation: dao.NewConversationDao(e),
		ConvMemb:     dao.NewConvMembDao(e),
	}
}
