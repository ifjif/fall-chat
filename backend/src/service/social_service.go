package service

import (
	"errors"
	"strconv"
	"time"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/utils"

	"gorm.io/gorm"
)

type SocialService struct {
	BaseService
	dao *data.SocialDao
}

func NewSocialService(name string) *SocialService {
	bs := BaseService{
		name: name,
	}

	return &SocialService{
		BaseService: bs,
	}
}

func (ss *SocialService) Inject() {
	ss.BaseService.Inject()
	ss.dao = data.NewSocialDao(ss.db)
}

func (ss *SocialService) ApplyAddFriend(cmd *command.SocialAddFriendCommand) error {
	targetId, err := ss.sf.Decode(cmd.TargetID)
	if err != nil {
		return ss.Wrap("解码目标id", err, "username", cmd.Username)
	}

	// 对方和我
	friendships, err := ss.dao.SearchFriendships(cmd.ID, targetId)
	if err != nil {
		return ss.Wrap("找好友关系", err, "username", cmd.Username)
	}

	for _, fs := range friendships {
		if fs.UserID == targetId && fs.FriendID == cmd.ID && fs.Status == 0 {
			return ss.Wrap("", errs.SocialErrBlacklisted, "username", cmd.Username)
		}

		if fs.UserID == cmd.ID && fs.FriendID == targetId && fs.Status == 0 {
			return ss.Wrap("", errs.SocialErrBlacklist, "username", cmd.Username)
		}
	}

	existsCount, err := ss.dao.SearchFriendRequestPendingAndAgreeCount(cmd.ID, targetId)
	if err != nil {
		return ss.Wrap("找好友申请数量", err, "username", cmd.Username)
	}

	if existsCount != 0 {
		return ss.Wrap("", errs.SocialErrFriendRequestExists, "username", cmd.Username)
	}

	// 添加申请信息
	fr := model.GenerateFriendRequest(cmd.ID, targetId, cmd.ApplyType, cmd.Greeting)
	err = utils.InjectID4Model(ss.sf, &fr.Model)
	if err != nil {
		return ss.Wrap("创建主键id", err, "username", cmd.Username)
	}

	err = ss.dao.CreateFriendRequest(fr)
	if err != nil {
		return ss.Wrap("创建好友申请", err, "username", cmd.Username)
	}

	return nil
}

// TODO:创建群申请，带有动态唯一键，这个错误没有具体捕获
// 找群和群申请创建 可放在一个事务中，如在创建群申请时，群解散或禁用了，这时可能出现脏数据
// (1: 数据库查添加行锁[太重，不太好] 2:采用redis锁[推荐])
func (ss *SocialService) ApplyAddGroup(cmd *command.SocialAddGroupRequestCommand) error {
	groupID, err := ss.sf.Decode(cmd.GroupID)
	if err != nil {
		return ss.Wrap("群ID解码", err, "username", cmd.Username)
	}

	userID := cmd.ID
	if cmd.InvitedID != "" {
		invitedID, err := ss.sf.Decode(cmd.InvitedID)
		if err != nil {
			return ss.Wrap("解码被邀请人ID", err, "username", cmd.Username)
		}

		if userID == invitedID {
			return ss.Wrap("", errs.SocialErrGroupRequestInviteSelf, "username", cmd.Username)
		}
		userID = invitedID

		// 邀请人需要在群
		member, err := ss.dao.FindGroupMemberWithRole(groupID, cmd.ID)
		if err != nil {
			return ss.Wrap("找群成员", err, "username", cmd.Username)
		}

		if member == nil {
			return ss.Wrap("", errs.SocialErrGroupMemberInvitorNotExists, "username", cmd.Username)
		}
	}

	existsCount, err := ss.dao.FindGroupRequestPendingAndAgreeCount(groupID, userID)
	if err != nil {
		return ss.Wrap("找群申请数量", err, "username", cmd.Username)
	}

	if existsCount > 0 {
		return ss.Wrap("", errs.SocialErrGroupRequestExists, "username", cmd.Username)
	}

	ctxFields := []any{
		"group_id", groupID,
		"username", cmd.Username,
	}
	groupInfo, err := ss.dao.FindGroupWithStatus(groupID)
	if err != nil {
		return ss.Wrap("找群", err, ctxFields...)
	}
	if groupInfo == nil {
		return ss.Wrap("", errs.SocialErrGroupNotExists, ctxFields...)
	}
	if groupInfo.IsDisabled() {
		return ss.Wrap("", errs.SocialErrGroupDisabled, ctxFields...)
	}
	if groupInfo.IsNotJoin() {
		return ss.Wrap("", errs.SocialErrGroupNotJoin, ctxFields...)
	}

	existsCount, err = ss.dao.ExistsGroupMember(groupID, userID)
	if err != nil {
		return ss.Wrap("找群成员数量", err, ctxFields...)
	}

	if existsCount > 0 {
		return ss.Wrap("", errs.SocialErrGroupMemberExists, ctxFields...)
	}

	// 添加申请
	gr := model.GroupRequest{}
	err = utils.InjectID4Model(ss.sf, &gr.Model)
	if err != nil {
		return ss.Wrap("生成群申请ID", err, ctxFields...)
	}
	gr.ToGroupID = groupID
	gr.FromUserID = userID
	if cmd.InvitedID != "" {
		gr.InviterUserID = cmd.ID
	}
	gr.Message = cmd.Greeting
	gr.ApplyType = int8(cmd.ApplyType)
	gr.ActiveUniqueKey()
	err = ss.dao.CreateOneGroupRequest(&gr)
	if err != nil {
		return ss.Wrap("创建群申请", err, ctxFields...)
	}
	return nil
}

func (ss *SocialService) SearchFriendRequest(qu *query.SocialSearchFriendRequestQuery) ([]*response.SocialSearchFriendRequestResult, error) {
	result := []*response.SocialSearchFriendRequestResult{}
	mo, err := ss.dao.SearchFriendRequestRow(qu.ID, qu.FilterStatus, 10)
	if err != nil {
		return nil, ss.Wrap("找好友请求信息", err, "id", qu.ID, "username", qu.Username)
	}

	for _, row := range mo {
		res := converter.ToSocialSearchFriendRequestResult(&row)
		res.Sender = row.FromUserID == qu.ID
		result = append(result, res)
	}

	return result, nil
}

func (ss *SocialService) HandleFriendRequest(cmd *command.SocialHandleFriendRequestCommand) error {
	targetID, err := ss.sf.Decode(cmd.HandleID)
	if err != nil {
		return ss.Wrap("id解码", err, "handle_id", cmd.HandleID, "username", cmd.Username)
	}

	req, err := ss.dao.SearchFriendshipByToUerID(targetID, cmd.ID)
	if err != nil {
		contextInfo := []any{"username", cmd.Username}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.SocialErrFriendRequestNotExists
		} else {
			contextInfo = append(contextInfo, "reason", "找好友请求任务失败")
		}
		return ss.Wrap("处理好友请求任务", err, contextInfo...)
	}

	// 必须是 待处理的
	if req.Status != 0 {
		return ss.Wrap("", errs.SocialErrFriendRequestHandledOrExpired, "username", cmd.Username)
	}

	var reason string
	err = ss.db.Transaction(func(tx *gorm.DB) error {
		reason = "更新好友请求状态失败"
		// 1同意，直接改状态，2:拒绝，该状态和uniqueKey设置timestamp
		nowStr := strconv.FormatInt(time.Now().UnixMilli(), 10)

		updates := map[string]any{
			"status": cmd.Action,
		}
		// 拒绝
		if cmd.Action != 1 {
			updates["unique_key"] = nowStr
		}

		err = ss.dao.UpdateFriendshipRequestStatusWithTx(tx, targetID, updates)
		if err != nil {
			return err
		}

		if cmd.Action != 1 {
			return nil
		}

		reason = "生成双向好友关系失败"

		// 同意
		aFriend := model.GenerateFriendship(cmd.ID, req.FromUserID, req.GetApplyTypeText())
		err = utils.InjectID4Model(ss.sf, &aFriend.Model)
		if err != nil {
			return err
		}

		bFriend := model.GenerateFriendship(req.FromUserID, cmd.ID, req.GetApplyTypeText())
		err = utils.InjectID4Model(ss.sf, &bFriend.Model)
		if err != nil {
			return err
		}

		reason = "保存双向好友关系到数据库失败"
		err = ss.dao.CreateFriendshipsWithTx(tx, []model.Friendship{*aFriend, *bFriend})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return ss.Wrap(reason, err, "username", cmd.Username, "action", cmd.Action)
	}

	return nil
}

func (ss *SocialService) FriendshipList(qu *query.SocialFriendshipListQuery) ([]model.FriendshipRow, error) {
	rows, err := ss.dao.SearchFriendshipList(qu.ID)
	if err != nil {
		return nil, ss.Wrap("查找好友列表", err, "username", qu.Username)
	}
	return rows, nil
}

func (ss *SocialService) SearchGroupRequest(qu *query.SocialSearchGroupRequestQuery) ([]*response.SocialSearchGroupRequestResult, error) {
	rows, err := ss.dao.SearchGroupRequestRow(qu.ID, qu.FilterStatus)
	if err != nil {
		return nil, ss.Wrap("找群申请信息", err, "username", qu.Username)
	}

	var result []*response.SocialSearchGroupRequestResult

	for i := range rows {
		v := converter.ToSocialSearchGroupRequestResult(&rows[i])
		v.Sender = rows[i].ApplyUser.ID == uint64(qu.ID)
		result = append(result, v)
	}
	return result, nil
}

func (ss *SocialService) HandleGroupRequest(cmd *command.SocialHandleGroupRequestCommand) error {
	// TODO: 注意并发处理，可加redis锁 , 群申请处理锁
	handleID, err := ss.sf.Decode(cmd.HandleID)
	if err != nil {
		return ss.Wrap("解码处理id", err, "handle_id", cmd.HandleID, "username", cmd.Username)
	}
	requestInfo, err := ss.dao.FindGroupRequestPendingInfo(handleID)
	if err != nil {
		return ss.Wrap("找待处理的群申请", err, "username", cmd.Username)
	}

	if requestInfo == nil {
		return ss.Wrap("", errs.SocialErrGroupRequestNotExists, "username", cmd.Username, "handleID", handleID)
	}

	if !requestInfo.IsPendingStatus() {
		return ss.Wrap("", errs.SocialErrGroupRequestHandledOrExpired, "username", cmd.Username,
			"handleID", handleID)
	}

	if requestInfo.FromUserID == cmd.HandleByUserID {
		return ss.Wrap("", errs.SocialErrGroupRequestHandleSelf, "username", cmd.Username,
			"HandleID", handleID)
	}

	// TODO: 群成员增加锁
	groupInfo, err := ss.dao.FindGroupWithMemberCount(requestInfo.ToGroupID)
	if err != nil {
		return ss.Wrap("找群信息", err, "username", cmd.Username)
	}
	if groupInfo.IsDisabled() {
		return ss.Wrap("", errs.SocialErrGroupDisabled, "username", cmd.Username)
	}

	// 判断成员数
	tempRequset := &model.GroupRequest{Status: int8(cmd.Action)}
	if tempRequset.IsApprovedStatus() && groupInfo.OverflowMember(1) {
		return ss.Wrap("", errs.SocialErrGroupMemberOverflow, "username", cmd.Username)
	}

	member, err := ss.dao.FindGroupMemberWithRole(requestInfo.ToGroupID, cmd.HandleByUserID)
	if err != nil {
		return ss.Wrap("找群成员角色", err, "username", cmd.Username)
	}

	if member == nil {
		return ss.Wrap("", errs.SocialErrGroupMemberNotExists, "username", cmd.Username)
	}
	if member.IsOrdinaryRole() {
		return ss.Wrap("", errs.SocialErrGroupMemberRoleOrdinary, "username", cmd.Username)
	}

	// 修改请求状态，如是同意，添加群成员
	reason := ""
	err = ss.db.Transaction(func(tx *gorm.DB) error {
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

		err = ss.dao.UpdateGroupRequestStatusWithTx(tx, updateRequst)
		if err != nil {
			return err
		}

		if !isApproved {
			return nil
		}

		reason = "生成群成员ID失败"
		member := &model.GroupMember{}
		err = utils.InjectID4Model(ss.sf, &member.Model)
		if err != nil {
			return err
		}
		member.GroupID = requestInfo.ToGroupID
		member.UserID = requestInfo.FromUserID
		member.JoinSource = requestInfo.GetApplyTypeText()

		reason = "创建群成员失败"
		err = ss.dao.CreateGroupMemberWithTx(tx, member)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return ss.Wrap(reason, err, "username", cmd.Username)
	}
	return nil
}
