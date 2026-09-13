package service

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/utils"

	"gorm.io/gorm"
)

type GroupService struct {
	BaseService
	dao *data.GroupDao
}

func NewGroupService(name string) *GroupService {
	return &GroupService{
		BaseService: BaseService{
			name: name,
		},
	}
}

func (gs *GroupService) Inject() {
	gs.BaseService.Inject()
	gs.dao = data.NewGroupDao(gs.db)
}

func (gs *GroupService) CreateGroup(cmd *command.CreateGroupCommand) error {
	groupInfo := model.GroupInfo{}
	err := utils.InjectID4Model(gs.sf, &groupInfo.Model)
	if err != nil {
		return gs.LogWarn("CreateGroup", err,
			"username", cmd.Username,
			"reason", "生成群ID失败")
	}

	searchFlag := int8(cmd.SearchFlag)
	joinMode := int8(cmd.JoinMode)
	groupInfo.GroupName = cmd.Name
	groupInfo.CreatorID = cmd.CreatorID
	groupInfo.MaxMembers = cmd.MaxMembers
	groupInfo.SearchFlag = &searchFlag
	groupInfo.JoinMode = &joinMode
	groupInfo.Description = cmd.Description
	groupInfo.Avatar = cmd.Avatar

	// 群主
	groupMember := model.GenerateCreatorGroupMember(int64(groupInfo.ID), cmd.CreatorID)
	err = utils.InjectID4Model(gs.sf, &groupMember.Model)
	if err != nil {
		return gs.LogWarn("CreateGroup", err, "username", cmd.Username,
			"reason", "生成群成员ID失败")
	}

	// 群会话
	groupConv := model.GenerateGroupConv(groupInfo.PublicID)
	err = utils.InjectID4Model(gs.sf, &groupConv.Model)
	if err != nil {
		return gs.LogWarn("CreateGroup", err, "username", cmd.Username,
			"reason", "生成群会话ID失败")
	}

	// 群会话成员
	groupConvMember := model.GenerateConvMember(int64(groupConv.ID), cmd.CreatorID, int64(groupInfo.ID), groupInfo.Avatar, groupInfo.GroupName)
	err = utils.InjectID4Model(gs.sf, &groupConvMember.Model)
	if err != nil {
		return gs.LogWarn("CreateGroup", err, "username", cmd.Username,
			"reason", "生成群会话成员ID失败")
	}

	reason := ""
	err = gs.db.Transaction(func(tx *gorm.DB) error {
		reason = "创建群失败"
		err = gs.dao.CreateOneGroupWithTx(tx, &groupInfo)
		if err != nil {
			return err
		}

		reason = "创建群主失败"
		err = gs.dao.CreateOneGroupMemberWithTx(tx, groupMember)
		if err != nil {
			return err
		}

		reason = "创建群会话失败"
		err = gs.dao.CreateOneGropConvWithTx(tx, groupConv)
		if err != nil {
			return err
		}

		reason = "创建群会话成员失败"
		err = gs.dao.CreateOneGroupConvMemberWithTx(tx, groupConvMember)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return gs.LogWarn("CreateGroup", err,
			"username", cmd.Username, "reason", reason)
	}
	return nil
}

func (gs *GroupService) Groups(qu *query.GroupInfoGroupsQuery) ([]*response.GroupInfoJoinedGroupsResult, error) {
	rows, err := gs.dao.FindJoinedGroups(qu.ID, qu.Filter)
	if err != nil {
		return nil, gs.LogWarn("Groups", err, "找我加入的群失败")
	}

	var result []*response.GroupInfoJoinedGroupsResult

	for i := range rows {
		result = append(result, converter.ToGroupInfoJoinedGroupsResult(&rows[i]))
	}

	return result, err
}
