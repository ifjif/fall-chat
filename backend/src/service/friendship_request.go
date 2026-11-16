package service

import (
	"context"
	"errors"
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
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/queryparam"
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
	targetId := cmd.TargetID
	userID := cmd.UserID
	// 对方和我
	friendships, err := frs.store.
		FriendshipRequestRepos().
		Friendship.
		FindRelationship(ctx, userID, targetId)
	if err != nil {
		return frs.Wrap("找好友关系", err, "username", cmd.Username)
	}

	for _, fs := range friendships {
		if fs.UserID == targetId && fs.FriendID == userID && fs.Status == 0 {
			return frs.Wrap("", errs.SocialErrBlacklisted, "username", cmd.Username)
		}

		if fs.UserID == userID && fs.FriendID == targetId && fs.Status == 0 {
			return frs.Wrap("", errs.SocialErrBlacklist, "username", cmd.Username)
		}
	}

	existsCount, err := frs.store.
		FriendshipRequestRepos().
		FriendshipRequest.
		PendingAndAgreeCount(ctx, userID, targetId)
	if err != nil {
		return frs.Wrap("找好友申请数量", err, "username", cmd.Username)
	}

	if existsCount != 0 {
		return frs.Wrap("", errs.SocialErrFriendRequestExists, "username", cmd.Username)
	}

	// 添加申请信息
	fr := model.GenerateFriendRequest(userID, targetId, cmd.ApplyType, cmd.Greeting)
	err = model.InjectID(&fr.Model)
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

func (frs *FriendshipRequest) List(ctx context.Context, qu *query.ListFriendRequest) (*response.PageCursor, error) {
	param := queryparam.FriendshipRequestList{
		Filter: filter.ListFilterOptions{
			FilterDef: filter.FriendshipRequestFilterDefinition,
			Filter:    filter.ToMap(qu.Filter),
		},
	}
	pager := &cursor.CursorPager[*cursor.FriendshipRequestCursor, *model.FriendshipRequestRow, *response.FriendRequestResult]{
		Limit: 20,
		DecodeNext: func(s string) (*cursor.FriendshipRequestCursor, error) {
			return cursor.Decode[cursor.FriendshipRequestCursor](s)
		},
		GetSortDef: func(key string) (cursor.SortDefinition, error) {
			sortDef, ok := cursor.FriendshipRequestSortDefinitions[key]
			if !ok {
				return cursor.EmptySortDef, errs.CommonErrSortKeyNotExist
			}
			return sortDef, nil
		},
		Query: func(ctx context.Context, sortValue map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.FriendshipRequestRow, error) {
			param.Sort = cursor.ListSortOptions{
				SortDef: sortDef,
				Sort:    sortValue,
			}
			return frs.store.
				FriendshipRequestRepos().
				FriendshipRequest.
				List(ctx, qu.UserID, param, limit)
		},
		BuildCursor: func(next *cursor.FriendshipRequestCursor, rows []*model.FriendshipRequestRow, lastIdx int, hasNext bool) *cursor.FriendshipRequestCursor {
			last := rows[lastIdx]
			c := &cursor.FriendshipRequestCursor{
				ID: last.ID,
			}
			c.SetHasNext(hasNext)
			return c
		},
		ToResult: func(fr *model.FriendshipRequestRow) *response.FriendRequestResult {
			return converter.ToFriendRequestResult(fr, int64(qu.UserID))
		},
	}

	result, encoded, hasNext, err := pager.List(ctx, qu.Next, qu.SortKey)
	if err != nil {
		return nil, frs.Wrap("list friendship request", err, "user_id", qu.UserID, "sort_key", qu.SortKey)
	}

	return &response.PageCursor{
		Data:    utils.ToAnySlice(result),
		Next:    encoded,
		HasNext: hasNext,
	}, nil
}

func (frs *FriendshipRequest) Handle(ctx context.Context, cmd *command.HandleFriendRequest) error {
	targetID := cmd.HandleID
	userID := cmd.UserID

	req, err := frs.store.
		FriendshipRequestRepos().
		FriendshipRequest.
		FindWithToUerID(ctx, targetID, userID)
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
		aFriend := model.GenerateFriendship(userID, req.FromUserID, req.GetApplyTypeText())
		err = model.InjectID(&aFriend.Model)
		if err != nil {
			return err
		}

		bFriend := model.GenerateFriendship(req.FromUserID, userID, req.GetApplyTypeText())
		err = model.InjectID(&bFriend.Model)
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
