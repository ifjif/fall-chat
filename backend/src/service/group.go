package service

import (
	"context"
	"encoding/json"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/repos"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/mess"
	messagequeue "zzc/fall-chat/src/message_queue"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/utils"

	"github.com/segmentio/kafka-go"
)

type Group struct {
	BaseService
	store store.GroupStore
}

func NewGroup(name string) *Group {
	return &Group{
		BaseService: BaseService{
			name: name,
		},
	}
}

func (gs *Group) Inject() {
	gs.BaseService.Inject()
	gs.store = data.GetBaseDao()
}

func (gs *Group) CreateGroup(ctx context.Context, cmd *command.CreateGroup) error {
	groupInfo := model.GroupInfo{}
	err := utils.InjectID4Model(gs.sf, &groupInfo.Model)
	if err != nil {
		return gs.Wrap("生成群ID", err, "id", cmd.CreatorID, "username", cmd.Username)
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
		return gs.Wrap("生成群成员ID", err, "username", cmd.Username)
	}

	// 群会话
	groupConv := model.GenerateGroupConv(groupInfo.PublicID)
	err = utils.InjectID4Model(gs.sf, &groupConv.Model)
	if err != nil {
		return gs.Wrap("生成群会话ID", err, "username", cmd.Username)
	}

	// 群会话成员
	groupConvMember := model.GenerateConvMember(int64(groupConv.ID), cmd.CreatorID, int64(groupInfo.ID))
	err = utils.InjectID4Model(gs.sf, &groupConvMember.Model)
	if err != nil {
		return gs.Wrap("生成群会话成员ID", err, "username", cmd.Username)
	}

	reason := ""
	gs.store.WithinGroupTX(ctx, func(repo repos.GroupRepositories) error {
		reason = "创建群失败"
		err = repo.Group.CreateOneGroup(ctx, &groupInfo)
		if err != nil {
			return err
		}

		reason = "创建群主失败"
		err = repo.GroupMember.CreateOneGroupMember(ctx, groupMember)
		if err != nil {
			return err
		}

		reason = "创建群会话失败"
		err = repo.Conversation.CreateOne(ctx, groupConv)
		if err != nil {
			return err
		}

		reason = "创建群会话成员失败"
		err = repo.ConvMemb.CreateOne(ctx, groupConvMember)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return gs.Wrap(reason, err, "username", cmd.Username)
	}

	// TODO: 发送消息队列, join_group
	v := mess.JoinGroup{
		GroupID: groupInfo.PublicID,
		UserID:  cmd.CreatorPublicID,
	}
	bv, _ := json.Marshal(v)
	go messagequeue.KWriter().WriteMessages(context.Background(), kafka.Message{
		Topic: "join-group",
		Key:   []byte(groupInfo.PublicID),
		Value: bv,
	})
	return nil
}

func (gs *Group) Groups(ctx context.Context, qu *query.GroupInfoGroups) ([]*response.GroupInfoJoinedGroupsResult, error) {
	rows, err := gs.store.
		GroupRepos().
		Group.
		FindJoinedGroups(ctx, qu.ID, qu.Filter)
	if err != nil {
		return nil, gs.Wrap("找我加入的群", err)
	}

	var result []*response.GroupInfoJoinedGroupsResult

	for i := range rows {
		result = append(result, converter.ToGroupInfoJoinedGroupsResult(&rows[i]))
	}

	return result, err
}
