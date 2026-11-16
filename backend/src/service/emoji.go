package service

import (
	"context"
	"fmt"
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
)

type Emoji struct {
	BaseService
	store store.EmojiStore
}

func NewEmoji(name string) *Emoji {
	bs := BaseService{name: name}
	return &Emoji{
		BaseService: bs,
	}
}

func (es *Emoji) Inject() {
	es.BaseService.Inject()
	es.store = data.GetBaseDao()
}

func (es *Emoji) CreatePack(ctx context.Context, cmd *command.CreateEmojiPack) error {
	m := model.EmojiPack{}
	err := model.InjectID(&m.Model)
	if err != nil {
		return es.Wrap("id生成", err)
	}
	m.Name = cmd.Name
	m.Source = cmd.Source
	m.Price = cmd.Price
	m.PriceType = cmd.PriceType
	m.CoverURL = cmd.CoverURL
	m.OwnerID = cmd.UserID
	m.CreatedAt = time.Now()

	err = es.store.
		EmojiRepos().
		EmojiPack.
		CreateOne(ctx, &m)
	if err != nil {
		return es.Wrap("创建EmojiPack", err)
	}

	return nil
}

func (es *Emoji) PackDetail(ctx context.Context, qu *query.EmojiPackDetail) (*response.EmojiPackResult, error) {
	packID := qu.PackageID

	row, err := es.store.
		EmojiRepos().
		EmojiPack.
		FindOne(ctx, packID)
	if err != nil {
		return nil, err
	}

	return converter.ToEmojiPackResult(row), nil
}

func (es *Emoji) ListPack(ctx context.Context, qu *query.ListEmojiPack) (*response.PageCursor, error) {
	param := queryparam.EmojiPackList{
		Filter: filter.ListFilterOptions{
			FilterDef: filter.EmojiPackFilterDefinition,
			Filter:    filter.ToMap(qu.Filter),
		},
	}

	pager := &cursor.CursorPager[*cursor.EmojiPackCurosr, *model.EmojiPack, *response.EmojiPackResult]{
		Limit: 20,
		DecodeNext: func(s string) (*cursor.EmojiPackCurosr, error) {
			return cursor.Decode[cursor.EmojiPackCurosr](s)
		},
		GetSortDef: func(key string) (cursor.SortDefinition, error) {
			def, ok := cursor.EmojiPackSortDefinitions[key]
			if !ok {
				return cursor.EmptySortDef, errs.CommonErrSortKeyNotExist
			}
			return def, nil
		},
		Query: func(ctx context.Context, sortValue map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.EmojiPack, error) {
			param.Sort = cursor.ListSortOptions{
				SortDef: sortDef,
				Sort:    sortValue,
			}
			return es.store.EmojiRepos().EmojiPack.List(ctx, param, limit)
		},
		BuildCursor: func(_ *cursor.EmojiPackCurosr, rows []*model.EmojiPack, lastIdx int, hasNext bool) *cursor.EmojiPackCurosr {
			last := rows[lastIdx]
			c := cursor.EmojiPackCurosr{
				ID: last.ID,
			}
			c.SetHasNext(hasNext)
			return &c
		},
		ToResult: converter.ToEmojiPackResult,
	}

	result, encoded, hasNext, err := pager.List(ctx, qu.Next, qu.SortKey)
	if err != nil {
		return nil, es.Wrap("list emoji_pack", err, "sort_key", qu.SortKey)
	}

	return &response.PageCursor{
		Data:    utils.ToAnySlice(result),
		Next:    *encoded,
		HasNext: hasNext,
	}, nil
}

func (es *Emoji) Publish(ctx context.Context, cmd *command.PublishEmoji) error {
	packID := cmd.PackID

	expectedItemCount := es.store.
		EmojiCache().
		GetPublishItemCount(ctx, cmd.UserID, uint64(packID))

	if expectedItemCount == 0 || expectedItemCount != len(cmd.Items) {
		return errs.EmojiErrPublishItemCountNotMatch
	}

	exist, err := es.store.EmojiRepos().
		EmojiPack.
		ExistsStatus(ctx, uint64(packID), cmd.UserID, model.NilEmojiPack.GetDraftStatusNilSafe())
	commonErrFields := []any{
		"user_id", cmd.UserID,
		"emoji_pack_id", packID,
		"emoji_status", 1,
	}
	if err != nil {
		return es.Wrap("找表情包状态存在", err, commonErrFields...)
	}

	if !exist {
		return es.Wrap("", errs.EmojiErrPackNotExists, commonErrFields...)
	}

	emojiItems := make([]model.EmojiItem, 0, len(cmd.Items))
	emojiOwns := make([]model.EmojiOwn, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		emojiItem := model.EmojiItem{}
		if err := model.InjectID(&emojiItem.Model); err != nil {
			return es.Wrap("emojiItem id生成失败", err, commonErrFields...)
		}
		emojiItem.PackID = uint64(packID)
		emojiItem.MediaType = item.MediaType
		emojiItem.MediaURL = item.MediaURL
		emojiItem.FontName = item.FontName
		emojiItem.ThumbURL = item.ThumbURL
		emojiItem.Width = item.Width
		emojiItem.Height = item.Height
		emojiItem.FileSize = item.FileSize
		emojiItem.SortNO = item.SortNO
		emojiItems = append(emojiItems, emojiItem)

		emojiOwn := model.EmojiOwn{}
		if err := model.InjectID(&emojiOwn.Model); err != nil {
			return es.Wrap("emojiOwn id生成失败", err, commonErrFields...)
		}
		emojiOwn.ItemID = emojiItem.ID
		emojiOwn.PackID = uint64(packID)
		emojiOwn.SetSelfSource()
		emojiOwn.UserID = cmd.UserID
		emojiOwns = append(emojiOwns, emojiOwn)
	}

	reason := ""
	err = es.store.WithinEmojiTX(ctx, func(repo repos.EmojiRepositories) error {
		reason = "更新表情包状态失败"
		err = repo.EmojiPack.
			UpdateStatus(ctx, uint64(packID), model.NilEmojiPack.GetAddedStatusNilSafe())
		if err != nil {
			return err
		}

		reason = "插入表情包item失败"
		err = repo.EmojiItem.CreateBatch(ctx, emojiItems, 200)
		if err != nil {
			return err
		}

		reason = "插入拥有自己上传表情信息失败"
		err = repo.EmojiOwn.CreateBatch(ctx, emojiOwns, 200)
		return err
	})
	if err != nil {
		commonErrFields = append(commonErrFields, "reason", reason)
		return es.Wrap("更新表情包状态，插入items", err, commonErrFields...)
	}

	return nil
}

// 校验并生成 临时推送token
func (es *Emoji) PublishVerifyCount(ctx context.Context, qu *query.PublishEmojiVerifyCount) error {
	maxItemsCount := 50
	if qu.ItemCount == 0 || qu.ItemCount > maxItemsCount {
		return errs.EmojiErrPublishOutOfBound.ReMessage(fmt.Sprintf("最多%d", maxItemsCount), true)
	}

	// 是否是自己的
	exists, err := es.store.
		EmojiRepos().
		EmojiPack.
		ExistsStatus(ctx, uint64(qu.PackID), qu.UserID, model.NilEmojiPack.GetDraftStatusNilSafe())
	if err != nil {
		return err
	}

	if !exists {
		return errs.EmojiErrPackNotExists
	}

	// 将 packId:userID:上传数量，存入redis
	err = es.store.
		EmojiCache().
		SetPublishItemCount(ctx, qu.UserID, uint64(qu.PackID), qu.ItemCount)
	if err != nil {
		return err
	}

	return nil
}

func (es *Emoji) ListItem(ctx context.Context, qu *query.ListEmojiItem) (*response.PageCursor, error) {
	pager := &cursor.CursorPager[*cursor.EmojiItemCursor, *model.EmojiItem, *response.EmojiItemResult]{
		Limit: 50,
		DecodeNext: func(s string) (*cursor.EmojiItemCursor, error) {
			return cursor.Decode[cursor.EmojiItemCursor](s)
		},
		GetSortDef: func(key string) (cursor.SortDefinition, error) {
			def, ok := cursor.EmojiItemSortDefinitions[key]
			if !ok {
				return cursor.EmptySortDef, errs.CommonErrSortKeyNotExist
			}
			return def, nil
		},
		Query: func(ctx context.Context, sortValue map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.EmojiItem, error) {
			return es.store.
				EmojiRepos().
				EmojiItem.
				List(ctx, uint64(qu.PackID), sortValue, sortDef, limit)
		},
		BuildCursor: func(next *cursor.EmojiItemCursor, rows []*model.EmojiItem, lastIdx int, hasNext bool) *cursor.EmojiItemCursor {
			last := rows[lastIdx]
			c := &cursor.EmojiItemCursor{
				ID:     last.ID,
				SortNo: last.SortNO,
			}
			c.SetHasNext(hasNext)
			return c
		},
		ToResult: converter.ToEmojiItemResult,
	}

	result, encoded, hasNext, err := pager.List(ctx, qu.Next, qu.SortKey)
	if err != nil {
		return nil, es.Wrap("list emoji item", err, "pack_id", qu.PackID, "sort_key", qu.SortKey)
	}

	return &response.PageCursor{
		Data:    utils.ToAnySlice(result),
		Next:    encoded,
		HasNext: hasNext,
	}, nil
}

func (es *Emoji) OwnListPack(ctx context.Context, qu *query.OwnListEmojiPack) (*response.PageCursor, error) {
	param := queryparam.EmojiOwnList{
		Owner: filter.ListFilterOptions{
			FilterDef: filter.EmojiOwnFilterDeinition,
			Filter:    filter.ToMap(qu.OwnFilter),
		},
		Final: filter.ListFilterOptions{
			FilterDef: filter.EmojiPackFilterDefinition,
			Filter:    filter.ToMap(qu.PackFilter),
		},
		SourceType: qu.OwnFilter.SourceType,
	}
	pager := &cursor.CursorPager[*cursor.EmojiPackCurosr, *model.EmojiPack, *response.EmojiPackResult]{
		Limit: 20,
		DecodeNext: func(s string) (*cursor.EmojiPackCurosr, error) {
			return cursor.Decode[cursor.EmojiPackCurosr](s)
		},
		GetSortDef: func(key string) (cursor.SortDefinition, error) {
			return cursor.EmojiPackSortDefinitions[key], nil
		},
		Query: func(ctx context.Context, sortValue map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.EmojiPack, error) {
			param.Sort = cursor.ListSortOptions{
				Sort:    sortValue,
				SortDef: sortDef,
			}

			return es.store.
				EmojiRepos().
				EmojiOwn.
				List(ctx, int64(qu.UserID), param, limit)
		},
		BuildCursor: func(next *cursor.EmojiPackCurosr, rows []*model.EmojiPack, lastIdx int, hasNext bool) *cursor.EmojiPackCurosr {
			last := rows[lastIdx]
			c := &cursor.EmojiPackCurosr{
				ID: last.ID,
			}

			c.SetHasNext(hasNext)
			return c
		},
		ToResult: converter.ToEmojiPackResult,
	}

	result, encoded, hasNext, err := pager.List(ctx, qu.Next, "created_desc")
	if err != nil {
		return nil, es.Wrap("list_own_emoji_pack", err, "user_id", qu.UserID)
	}
	return &response.PageCursor{
		Data:    utils.ToAnySlice(result),
		Next:    encoded,
		HasNext: hasNext,
	}, nil
}
