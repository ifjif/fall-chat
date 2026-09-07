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
		ss.LogWarn("ApplyAddFriend", err,
			"username", cmd.Username,
			"reason", "解码目标id失败")
		return err
	}

	// 对方和我
	friendships, err := ss.dao.SearchFriendships(cmd.ID, targetId)
	if err != nil {
		ss.LogWarn("ApplyAddFriend", err, "username", cmd.Username,
			"reason", "找好友关系错误",
		)
		return err
	}

	for _, fs := range friendships {
		if fs.UserID == targetId && fs.FriendID == cmd.ID && fs.Status == 0 {
			ss.LogWarn("ApplyAddFriend", errs.SocialErrBlacklisted, "username", cmd.Username)
			return errs.SocialErrBlacklist
		}

		if fs.UserID == cmd.ID && fs.FriendID == targetId && fs.Status == 0 {
			ss.LogWarn("ApplyAddFriend", errs.SocialErrBlacklist, "username", cmd.Username)
			return errs.SocialErrBlacklist
		}
	}

	existsCount, err := ss.dao.SearchFriendRequestPendingAndAgreeCount(cmd.ID, targetId)
	if err != nil {
		ss.LogWarn("ApplyAddFriend", err, "username", cmd.Username,
			"reason", "找好友申请数量错误",
		)

		return err
	}

	if existsCount != 0 {
		ss.LogWarn("ApplyAddFriend", errs.SocialErrFriendRequestExists, "username", cmd.Username)
		return errs.SocialErrFriendRequestExists
	}

	// 添加申请信息
	fr := model.GenerateFriendRequest(cmd.ID, targetId, cmd.ApplyType, cmd.Greeting)
	err = utils.InjectID4Model(ss.sf, &fr.Model)
	if err != nil {
		ss.LogWarn("ApplyAddFriend", err, "username", cmd.Username,
			"reason", "创建主键id失败",
		)

		return err
	}

	err = ss.dao.CreateFriendRequest(fr)
	if err != nil {
		ss.LogWarn("ApplyAddFriend", err, "username", cmd.Username,
			"reason", "创建好友申请失败",
		)

		return err
	}

	return nil
}

// TODO:创建群申请，带有动态唯一键，这个错误没有具体捕获
// 找群和群申请创建 可放在一个事务中，如在创建群申请时，群解散或禁用了，这时可能出现脏数据
// (1: 数据库查添加行锁[太重，不太好] 2:采用redis锁[推荐])
func (ss *SocialService) ApplyAddGroup(cmd *command.SocialAddGroupRequestCommand) error {
	groupID, err := ss.sf.Decode(cmd.GroupID)
	if err != nil {
		ss.LogWarn("ApplyAddGroup", err, "username", cmd.Username,
			"reason", "群ID解码失败")
		return err
	}

	userID := cmd.ID
	if cmd.InvitedID != "" {
		invitedID, err := ss.sf.Decode(cmd.InvitedID)
		if err != nil {
			ss.LogWarn("ApplyAddGroup", err, "username", cmd.Username,
				"reason", "解码被邀请人ID失败")
			return err
		}

		if userID == invitedID {
			ss.LogWarn("ApplyAddGroup", errs.SocialErrGroupRequestInviteSelf, "username", cmd.Username)
			return errs.SocialErrGroupRequestInviteSelf
		}
		userID = invitedID

		// 邀请人需要在群
		member, err := ss.dao.FindGroupMemberWithRole(groupID, cmd.ID)
		if err != nil {
			ss.LogWarn("ApplyAddGroup", err, "username", cmd.Username,
				"reason", "找群成员失败")
			return err
		}

		if member == nil {
			ss.LogWarn("ApplyAddGroup", errs.SocialErrGroupMemberInvitorNotExists, "username", cmd.Username)
			return errs.SocialErrGroupMemberInvitorNotExists
		}
	}

	existsCount, err := ss.dao.FindGroupRequestPendingAndAgreeCount(groupID, userID)
	if err != nil {
		ss.LogWarn("ApplyAddGroup", err, "username", cmd.Username,
			"reason", "找群申请数量失败")
		return err
	}

	if existsCount > 0 {
		ss.LogWarn("ApplyAddGroup", errs.SocialErrGroupRequestExists, "username", cmd.Username)
		return errs.SocialErrGroupRequestExists
	}

	groupInfo, err := ss.dao.FindGroupWithStatus(groupID)
	if err != nil {
		ss.LogWarn("ApplyAddGroup", err, "username", cmd.Username,
			"reason", "找群失败")
		return err
	}
	if groupInfo == nil {
		ss.LogWarn("ApplyAddGroup", errs.SocialErrGroupNotExists, "username", cmd.Username)
		return errs.SocialErrGroupNotExists
	}
	if groupInfo.IsDisabled() {
		ss.LogWarn("ApplyAddGroup", errs.SocialErrGroupDisabled, "username", cmd.Username)
		return errs.SocialErrGroupDisabled
	}
	if groupInfo.IsNotJoin() {
		ss.LogWarn("ApplyAddGroup", errs.SocialErrGroupNotJoin, "username", cmd.Username)
		return errs.SocialErrGroupNotJoin
	}

	existsCount, err = ss.dao.ExistsGroupMember(groupID, userID)
	if err != nil {
		ss.LogWarn("ApplyAddGroup", err, "username", cmd.Username,
			"reason", "找群成员数量失败")
		return err
	}

	if existsCount > 0 {
		ss.LogWarn("ApplyAddGroup", errs.SocialErrGroupMemberExists, "username", cmd.Username)
		return errs.SocialErrGroupMemberExists
	}

	// 添加申请
	gr := model.GroupRequest{}
	err = utils.InjectID4Model(ss.sf, &gr.Model)
	if err != nil {
		ss.LogWarn("ApplyAddGroup", err, "username", cmd.Username,
			"reason", "生成群申请ID失败")
		return err
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
		ss.LogWarn("ApplyAddGroup", err, "username", cmd.Username,
			"reason", "创建群申请失败")
		return err
	}
	return nil
}

func (ss *SocialService) SearchFriendRequest(qu *query.SocialSearchFriendRequestQuery) ([]*response.SocialSearchFriendRequestResult, error) {
	result := []*response.SocialSearchFriendRequestResult{}
	mo, err := ss.dao.SearchFriendRequestRow(qu.ID, qu.FilterStatus, 10)
	if err != nil {
		ss.LogWarn("SearchFriendRequest", err, "找好友请求信息失败", "username", qu.Username)
		return nil, err
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
		ss.LogWarn("HandleFriendRequest", err, "username", cmd.Username, "reason", "id解码失败")
		return err
	}

	req, err := ss.dao.SearchFriendshipByToUerID(targetID, cmd.ID)
	if err != nil {
		contextInfo := []any{"username", cmd.Username}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.SocialErrFriendRequestNotExists
		} else {
			contextInfo = append(contextInfo, "reason", "找好友请求任务失败")
		}
		ss.LogWarn("HandleFriendRequest", err, contextInfo...)
		return err
	}

	// 必须是 待处理的
	if req.Status != 0 {
		ss.LogWarn("HandleFriendRequest", errs.SocialErrFriendRequestHandledOrExpired, "username", cmd.Username)
		return errs.SocialErrFriendRequestHandledOrExpired
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
		ss.LogWarn("HandleFriendRequest", err, "username", cmd.Username,
			"action", cmd.Action,
			"reason", reason)
		return err
	}

	return nil
}

func (ss *SocialService) FriendshipList(qu *query.SocialFriendshipListQuery) ([]model.FriendshipRow, error) {
	rows, err := ss.dao.SearchFriendshipList(qu.ID)
	if err != nil {
		ss.LogWarn("FriendshipList", err, "username", qu.Username,
			"reason", "查找好友列表失败")
		return nil, err
	}
	return rows, nil
}

func (ss *SocialService) SearchGroupRequest(qu *query.SocialSearchGroupRequestQuery) ([]*response.SocialSearchGroupRequestResult, error) {
	rows, err := ss.dao.SearchGroupRequestRow(qu.ID, qu.FilterStatus)
	if err != nil {
		ss.LogWarn("SearchGroupRequest", err, "username", qu.Username,
			"reason", "找群申请信息失败")
		return nil, err
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
		ss.LogWarn("HandleGroupRequest", err, "username", cmd.Username,
			"reason", "解码处理id失败")
		return err
	}
	requestInfo, err := ss.dao.FindGroupRequestPendingInfo(handleID)
	if err != nil {
		ss.LogWarn("HandleGroupRequest", err, "username", cmd.Username,
			"reason", "找待处理的群申请失败")
		return err
	}

	if requestInfo == nil {
		ss.LogWarn("HandleGroupRequest", errs.SocialErrGroupRequestNotExists, "username", cmd.Username,
			"handleID", handleID)
		return errs.SocialErrGroupRequestNotExists
	}

	if !requestInfo.IsPendingStatus() {
		ss.LogWarn("HandleGroupRequest", errs.SocialErrGroupRequestHandledOrExpired, "username", cmd.Username,
			"handleID", handleID)
		return errs.SocialErrGroupRequestHandledOrExpired
	}

	if requestInfo.FromUserID == cmd.HandleByUserID {
		ss.LogWarn("HandleGroupRequest", errs.SocialErrGroupRequestHandleSelf, "username", cmd.Username,
			"HandleID", handleID)
		return errs.SocialErrGroupRequestHandleSelf
	}

	// TODO: 群成员增加锁
	groupInfo, err := ss.dao.FindGroupWithMemberCount(requestInfo.ToGroupID)
	if err != nil {
		ss.LogWarn("HandleGroupRequest", err, "username", cmd.Username,
			"reason", "找群信息失败")
		return err
	}
	if groupInfo.IsDisabled() {
		ss.LogWarn("HandleGroupRequest", errs.SocialErrGroupDisabled, "username", cmd.Username)
		return errs.SocialErrGroupDisabled
	}

	// 判断成员数
	tempRequset := &model.GroupRequest{Status: int8(cmd.Action)}
	if tempRequset.IsApprovedStatus() && groupInfo.OverflowMember(1) {
		ss.LogWarn("HandleGroupRequest", errs.SocialErrGroupMemberOverflow, "username", cmd.Username)
		return errs.SocialErrGroupMemberOverflow
	}

	member, err := ss.dao.FindGroupMemberWithRole(requestInfo.ToGroupID, cmd.HandleByUserID)
	if err != nil {
		ss.LogWarn("HandleGroupRequest", err, "username", cmd.Username,
			"reason", "找群成员角色失败")
		return err
	}

	if member == nil {
		ss.LogWarn("HandleGroupRequest", errs.SocialErrGroupMemberNotExists, "username", cmd.Username)
		return errs.SocialErrGroupMemberNotExists
	}
	if member.IsOrdinaryRole() {
		ss.LogWarn("HandleGroupRequest", errs.SocialErrGroupMemberRoleOrdinary, "username", cmd.Username)
		return errs.SocialErrGroupMemberRoleOrdinary
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
		ss.LogWarn("HandleGroupRequest", err, "username", cmd.Username,
			"reason", reason)
	}
	return nil
}
