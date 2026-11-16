package service

import (
	"context"
	"strconv"
	"time"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/cursor"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/repos"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/filter"
	"zzc/fall-chat/src/mess"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/mq/producer"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/queryparam"
	"zzc/fall-chat/src/utils"
)

type GroupRequest struct {
	BaseService
	store             store.GroupRequestStore
	joinGroupProducer *producer.JoinGroup
}

func NewGroupRequest(name string) *GroupRequest {
	bs := BaseService{name: name}
	return &GroupRequest{
		BaseService: bs,
	}
}

func (grs *GroupRequest) Inject() {
	grs.BaseService.Inject()
	grs.store = data.GetBaseDao()
	grs.joinGroupProducer = producer.GetJoinGroup()
}

// 创建入群申请时可能群没有了，不额外处理
// 留给处理 请求来处理
// 这里只负责 创建记录
func (grs *GroupRequest) Apply(ctx context.Context, cmd *command.AddGroupRequest) error {
	groupID := cmd.GroupID

	userID := cmd.UserID
	inviterID := cmd.UserID
	if cmd.InvitedID != nil {
		invitedID := *cmd.InvitedID

		if userID == invitedID {
			return grs.Wrap("", errs.SocialErrGroupRequestInviteSelf, "username", cmd.Username)
		}
		userID = invitedID

		// 邀请人需要在群
		member, err := grs.store.
			GroupRequestRepos().
			GroupMemeber.
			FindRole(ctx, groupID, inviterID)
		if err != nil {
			return grs.Wrap("找群成员", err, "username", cmd.Username)
		}

		if member == nil {
			return grs.Wrap("", errs.SocialErrGroupMemberInvitorNotExists, "username", cmd.Username)
		}
	}

	// 等待和同意的数量
	applyCount, err := grs.store.
		GroupRequestRepos().
		GroupRequest.
		PendingAndAgreeCount(ctx, groupID, userID)
	if err != nil {
		return grs.Wrap("找群申请数量", err, "username", cmd.Username)
	}

	if applyCount > 0 {
		return grs.Wrap("", errs.SocialErrGroupRequestExists, "username", cmd.Username)
	}

	// 群的 状态和 加入模式
	ctxFields := []any{
		"group_id", groupID,
		"username", cmd.Username,
	}
	groupInfo, err := grs.store.
		GroupRequestRepos().
		Group.
		FindStatusAndJoinMode(ctx, groupID)
	if err != nil {
		return grs.Wrap("找群", err, ctxFields...)
	}
	if groupInfo == nil {
		return grs.Wrap("", errs.SocialErrGroupNotExists, ctxFields...)
	}
	if groupInfo.IsDisabled() {
		return grs.Wrap("", errs.SocialErrGroupDisabled, ctxFields...)
	}
	if groupInfo.IsNotJoin() {
		return grs.Wrap("", errs.SocialErrGroupNotJoin, ctxFields...)
	}

	// 是否在群
	exists, err := grs.store.
		GroupRequestRepos().
		GroupMemeber.
		Exists(ctx, groupID, userID)
	if err != nil {
		return grs.Wrap("找群成员数量", err, ctxFields...)
	}

	if exists > 0 {
		return grs.Wrap("", errs.SocialErrGroupMemberExists, ctxFields...)
	}

	// 添加申请
	gr := model.GroupRequest{}
	err = model.InjectID(&gr.Model)
	if err != nil {
		return grs.Wrap("生成群申请ID", err, ctxFields...)
	}
	gr.ToGroupID = groupID
	gr.FromUserID = userID
	if cmd.InvitedID != nil {
		gr.InviterUserID = inviterID
	}
	gr.Message = cmd.Greeting
	gr.ApplyType = int8(cmd.ApplyType)
	gr.ActiveUniqueKey()
	err = grs.store.
		GroupRequestRepos().
		GroupRequest.
		CreateOne(ctx, &gr)
		// INFO: 可能是唯一键冲突，说明已经存在
	if err != nil {
		return grs.Wrap("创建入群申请", err, ctxFields...)
	}
	return nil
}

func (grs *GroupRequest) List(ctx context.Context, qu *query.ListGroupRequest) (*response.PageCursor, error) {
	param := queryparam.GroupRequestList{
		Filter: filter.ListFilterOptions{
			FilterDef: filter.GroupRequestFilterDefinition,
			Filter:    filter.ToMap(qu.Filter),
		},
	}
	pager := &cursor.CursorPager[*cursor.GroupRequestCursor, *model.GroupRequestRow, *response.GroupRequestResult]{
		Limit: 20,
		DecodeNext: func(s string) (*cursor.GroupRequestCursor, error) {
			return cursor.Decode[cursor.GroupRequestCursor](s)
		},
		GetSortDef: func(key string) (cursor.SortDefinition, error) {
			def, ok := cursor.GroupRequestSortDefinitions[key]
			if !ok {
				return cursor.EmptySortDef, errs.CommonErrSortKeyNotExist
			}
			return def, nil
		},
		Query: func(ctx context.Context, sortValue map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.GroupRequestRow, error) {
			param.Sort = cursor.ListSortOptions{
				SortDef: sortDef,
				Sort:    sortValue,
			}
			return grs.store.
				GroupRequestRepos().
				GroupRequest.
				List(ctx, uint64(qu.UserID), param, limit)
		},
		BuildCursor: func(next *cursor.GroupRequestCursor, rows []*model.GroupRequestRow, lastIdx int, hasNext bool) *cursor.GroupRequestCursor {
			last := rows[lastIdx]
			c := &cursor.GroupRequestCursor{
				ID: last.ID,
			}
			c.SetHasNext(hasNext)
			return c
		},
		ToResult: func(row *model.GroupRequestRow) *response.GroupRequestResult {
			return converter.ToGroupRequestResult(row, uint64(qu.UserID))
		},
	}

	result, encoded, hasNext, err := pager.List(ctx, qu.Next, qu.SortKey)
	if err != nil {
		return nil, grs.Wrap("list group_request", err, "user_id", qu.UserID, "sort_key", qu.SortKey)
	}

	return &response.PageCursor{
		Data:    utils.ToAnySlice(result),
		Next:    encoded,
		HasNext: hasNext,
	}, nil
}

//	申请 只能被一次处理
//
// 判断群成员 (可增加锁)
func (grs *GroupRequest) Handle(ctx context.Context, cmd *command.HandleGroupRequest) error {
	handleID := cmd.HandleID

	// 找群申请
	requestInfo, err := grs.store.
		GroupRequestRepos().
		GroupRequest.FindBriefOne(ctx, handleID)
	if err != nil {
		return grs.Wrap("找待处理的群申请", err, "username", cmd.Username)
	}

	if requestInfo == nil {
		return grs.Wrap("", errs.SocialErrGroupRequestNotExists, "username", cmd.Username, "handleID", handleID)
	}

	if !requestInfo.IsPendingStatus() {
		return grs.Wrap("", errs.SocialErrGroupRequestHandledOrExpired, "username", cmd.Username,
			"handleID", handleID)
	}

	if requestInfo.FromUserID == cmd.HandleByUserID {
		return grs.Wrap("", errs.SocialErrGroupRequestHandleSelf, "username", cmd.Username,
			"HandleID", handleID)
	}

	// 找群 成员数和 最大成员数
	groupInfo, err := grs.store.
		GroupRequestRepos().
		Group.
		FindStatusAndMembCountAndMaxMemb(ctx, requestInfo.ToGroupID)
	if err != nil {
		return grs.Wrap("找群信息", err, "username", cmd.Username)
	}
	if groupInfo.IsDisabled() {
		return grs.Wrap("", errs.SocialErrGroupDisabled, "username", cmd.Username)
	}

	// 判断成员数
	tempRequset := &model.GroupRequest{Status: int8(cmd.Action)}
	if tempRequset.IsApprovedStatus() && groupInfo.OverflowMember(1) {
		return grs.Wrap("", errs.SocialErrGroupMemberOverflow, "username", cmd.Username)
	}

	// 找处理人 角色
	member, err := grs.store.
		GroupRequestRepos().
		GroupMemeber.
		FindRole(ctx, requestInfo.ToGroupID, cmd.HandleByUserID)
	if err != nil {
		return grs.Wrap("找群成员角色", err, "username", cmd.Username)
	}

	if member == nil {
		return grs.Wrap("", errs.SocialErrGroupMemberNotExists, "username", cmd.Username)
	}
	if member.IsOrdinaryRole() {
		return grs.Wrap("", errs.SocialErrGroupMemberRoleOrdinary, "username", cmd.Username)
	}

	groupID := model.EncodeID(requestInfo.ToGroupID)

	// 找群会话
	convID, err := grs.store.
		GroupRequestRepos().
		Group.
		FindConvID(ctx, groupID)
	if err != nil {
		return grs.Wrap("找群会话ID", err, "single_key", groupID)
	}
	if convID == 0 {
		return grs.Wrap("", errs.SocialErrGroupConvNotExists, "single_key", groupID)
	}

	// 修改请求状态，如是同意，添加群成员
	reason := ""
	grs.store.WithinGroupRequestTX(ctx, func(repo repos.GroupRequestRepositories) error {
		reason = "更新群申请状态失败"
		updateRequst := &model.GroupRequest{}
		updateRequst.ID = uint64(handleID)
		updateRequst.Status = int8(cmd.Action)
		updateRequst.HandleByUserID = cmd.HandleByUserID
		isApproved := updateRequst.IsApprovedStatus()
		if !isApproved {
			timeStr := strconv.FormatInt(time.Now().UnixMilli(), 10)
			updateRequst.UniqueKey = timeStr
		}

		err = repo.GroupRequest.UpdateOne(ctx, updateRequst)
		if err != nil {
			return err
		}

		if !isApproved {
			return nil
		}

		reason = "生成群成员ID失败"
		member := &model.GroupMember{}
		err = model.InjectID(&member.Model)
		if err != nil {
			return err
		}
		member.GroupID = requestInfo.ToGroupID
		member.UserID = requestInfo.FromUserID
		member.JoinSource = requestInfo.GetApplyTypeText()

		reason = "创建群成员失败"
		err = repo.GroupMemeber.CreateOne(ctx, member)
		if err != nil {
			return err
		}

		reason = "生成群会话成员ID失败"
		convMember := model.GenerateConvMember(convID, member.UserID, member.GroupID)
		err = model.InjectID(&convMember.Model)
		if err != nil {
			return err
		}

		reason = "创建群会话成员失败"
		err = repo.ConvMemb.CreateOne(ctx, convMember)

		return nil
	})

	if err != nil {
		return grs.Wrap(reason, err, "username", cmd.Username)
	}

	userID := model.EncodeID(requestInfo.FromUserID)

	v := mess.JoinGroup{
		GroupID: groupID,
		UserID:  userID,
	}

	grs.joinGroupProducer.Publish(ctx, groupID, &v)
	return nil
}
