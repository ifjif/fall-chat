package service

import (
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/model"
)

type GroupService struct {
	BaseService
}

func NewGroupService(name string) *GroupService {
	return &GroupService{
		BaseService: BaseService{
			name: name,
		},
	}
}

func (gs *GroupService) CreateGroup(cmd *command.CreateGroupCommand) error {
	id, err := gs.sf.NextID()
	if err != nil {
		gs.LogWarn("CreateGroup", err,
			"username", cmd.Username,
			"reason", "生成主键失败")
		return err
	}

	publicID := gs.sf.Encode(id)

	groupInfo := model.GroupInfo{
		GroupName:   cmd.Name,
		CreatorID:   cmd.CreatorID,
		MaxMembers:  cmd.MaxMembers,
		SearchFlag:  int8(cmd.SearchFlag),
		JoinMode:    int8(cmd.JoinMode),
		Description: cmd.Description,
		Avatar:      cmd.Avatar,
	}

	groupInfo.ID = uint64(id)
	groupInfo.PublicID = publicID

	result := gs.db.Create(&groupInfo).Omit("status", "MutedAll")
	if result.Error != nil {
		gs.LogWarn("CreateGroup", result.Error,
			"username", cmd.Username)

		return result.Error
	}
	return nil
}
