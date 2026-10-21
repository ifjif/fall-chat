package service

import (
	"context"
	"errors"
	"strconv"
	"time"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/repos"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/utils"

	"gorm.io/gorm"
)

type FriendshipRequest struct {
	BaseService
	store store.FriendshipRequestStore
}

func NewFriendshipRequest(name string) *FriendshipRequest {
	bs := BaseService{name: name}
	return &FriendshipRequest{
		BaseService: bs,
	}
}

func (frs *FriendshipRequest) Inject() {
	frs.BaseService.Inject()
	frs.store = data.GetBaseDao()
}

func (frs *FriendshipRequest) Apply(ctx context.Context, cmd *command.AddFriend) error {
	targetId, err := frs.sf.Decode(cmd.TargetID)
	if err != nil {
		return frs.Wrap("解码目标id", err, "username", cmd.Username)
	}

	// 对方和我
	friendships, err := frs.store.
		FriendshipRequestRepos().
		Friendship.
		FindRelationship(ctx, cmd.ID, targetId)
	if err != nil {
		return frs.Wrap("找好友关系", err, "username", cmd.Username)
	}

	for _, fs := range friendships {
		if fs.UserID == targetId && fs.FriendID == cmd.ID && fs.Status == 0 {
			return frs.Wrap("", errs.SocialErrBlacklisted, "username", cmd.Username)
		}

		if fs.UserID == cmd.ID && fs.FriendID == targetId && fs.Status == 0 {
			return frs.Wrap("", errs.SocialErrBlacklist, "username", cmd.Username)
		}
	}

	existsCount, err := frs.store.
		FriendshipRequestRepos().
		FriendshipRequest.
		PendingAndAgreeCount(ctx, cmd.ID, targetId)
	if err != nil {
		return frs.Wrap("找好友申请数量", err, "username", cmd.Username)
	}

	if existsCount != 0 {
		return frs.Wrap("", errs.SocialErrFriendRequestExists, "username", cmd.Username)
	}

	// 添加申请信息
	fr := model.GenerateFriendRequest(cmd.ID, targetId, cmd.ApplyType, cmd.Greeting)
	err = utils.InjectID4Model(frs.sf, &fr.Model)
	if err != nil {
		return frs.Wrap("创建主键id", err, "username", cmd.Username)
	}

	err = frs.store.
		FriendshipRequestRepos().
		FriendshipRequest.
		CreateOne(ctx, fr)
	if err != nil {
		return frs.Wrap("创建好友申请", err, "username", cmd.Username)
	}

	return nil
}

func (frs *FriendshipRequest) List(ctx context.Context, qu *query.ListFriendRequest) ([]*response.ListFriendRequestResult, error) {
	result := []*response.ListFriendRequestResult{}
	mo, err := frs.store.
		FriendshipRequestRepos().
		FriendshipRequest.
		List(ctx, qu.ID, qu.FilterStatus, 10)
	if err != nil {
		return nil, frs.Wrap("找好友请求信息", err, "id", qu.ID, "username", qu.Username)
	}

	for _, row := range mo {
		res := converter.ToListFriendRequestResult(&row)
		res.Sender = row.FromUserID == qu.ID
		result = append(result, res)
	}

	return result, nil
}

func (frs *FriendshipRequest) Handle(ctx context.Context, cmd *command.HandleFriendRequest) error {
	targetID, err := frs.sf.Decode(cmd.HandleID)
	if err != nil {
		return frs.Wrap("id解码", err, "handle_id", cmd.HandleID, "username", cmd.Username)
	}

	req, err := frs.store.
		FriendshipRequestRepos().
		FriendshipRequest.
		FindWithToUerID(ctx, targetID, cmd.ID)
	if err != nil {
		contextInfo := []any{"username", cmd.Username}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.SocialErrFriendRequestNotExists
		} else {
			contextInfo = append(contextInfo, "reason", "找好友请求任务失败")
		}
		return frs.Wrap("处理好友请求任务", err, contextInfo...)
	}

	// 必须是 待处理的
	if req.Status != 0 {
		return frs.Wrap("", errs.SocialErrFriendRequestHandledOrExpired, "username", cmd.Username)
	}

	var reason string

	err = frs.store.WithinFriendshipRequestTX(ctx, func(repo repos.FriendshipRequestRepositories) error {
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

		err = repo.FriendshipRequest.UpdateStatus(ctx, targetID, updates)
		if err != nil {
			return err
		}

		if cmd.Action != 1 {
			return nil
		}

		reason = "生成双向好友关系失败"

		// 同意
		aFriend := model.GenerateFriendship(cmd.ID, req.FromUserID, req.GetApplyTypeText())
		err = utils.InjectID4Model(frs.sf, &aFriend.Model)
		if err != nil {
			return err
		}

		bFriend := model.GenerateFriendship(req.FromUserID, cmd.ID, req.GetApplyTypeText())
		err = utils.InjectID4Model(frs.sf, &bFriend.Model)
		if err != nil {
			return err
		}

		reason = "保存双向好友关系到数据库失败"
		err = repo.Friendship.Create(ctx, []model.Friendship{*aFriend, *bFriend})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return frs.Wrap(reason, err, "username", cmd.Username, "action", cmd.Action)
	}

	return nil
}
