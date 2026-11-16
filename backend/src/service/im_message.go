package service

import (
	"context"
	"zzc/fall-chat/src/converter"
	"zzc/fall-chat/src/cursor"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"
	"zzc/fall-chat/src/utils"
)

type IMMessage struct {
	BaseService
	store store.IMMessageStore
}

func NewIMMessage(name string) *IMMessage {
	bs := BaseService{name: name}
	return &IMMessage{
		BaseService: bs,
	}
}

func (imm *IMMessage) Inject() {
	imm.BaseService.Inject()
	imm.store = data.GetBaseDao()
}

func (imm *IMMessage) List(ctx context.Context, qu *query.ListIMMessage) (*response.PageCursor, error) {
	ctxFields := []any{
		"cursor", qu.Next,
		"sort_key", qu.SortKey,
		"conv_id", qu.ConvID,
		"user_id", qu.UserID,
		"username", qu.Username,
	}
	convID := qu.ConvID

	pager := &cursor.CursorPager[*cursor.IMMessageCursor, *model.IMMessage, *response.IMMessageResult]{
		Limit: 30,
		DecodeNext: func(s string) (*cursor.IMMessageCursor, error) {
			return cursor.Decode[cursor.IMMessageCursor](s)
		},
		GetSortDef: func(key string) (cursor.SortDefinition, error) {
			def, ok := cursor.IMMessageSortDefinitions[key]
			if !ok {
				return cursor.EmptySortDef, errs.CommonErrSortKeyNotExist
			}
			return def, nil
		},
		Query: func(ctx context.Context, sortValue map[string]any, sortDef cursor.SortDefinition, limit int) ([]*model.IMMessage, error) {
			return imm.store.
				IMMessageRepos().
				IMMessage.
				List(ctx, uint64(convID), sortValue, sortDef, limit)
		},
		BuildCursor: func(next *cursor.IMMessageCursor, rows []*model.IMMessage, lastIdx int, hasNext bool) *cursor.IMMessageCursor {
			last := rows[lastIdx]
			c := &cursor.IMMessageCursor{
				Seq: last.Seq,
				ID:  last.ID,
			}
			c.SetHasNext(hasNext)
			return c
		},
		ToResult: converter.ToImMessageResult,
	}

	result, encoded, hasNext, err := pager.List(ctx, qu.Next, qu.SortKey)
	if err != nil {
		return nil, imm.Wrap("list im_message", err, ctxFields...)
	}

	return &response.PageCursor{
		Data:    utils.ToAnySlice(result),
		Next:    encoded,
		HasNext: hasNext,
	}, nil
}

func (imm *IMMessage) PatchList(ctx context.Context, qu *query.PatchListIMMessage) ([]*response.IMMessageResult, error) {
	convID := qu.ConvID

	rows, err := imm.store.
		IMMessageRepos().
		IMMessage.
		PatchList(ctx, uint64(convID), qu.MinSeq, qu.MaxSeq)
	if err != nil {
		return nil, err
	}

	result := make([]*response.IMMessageResult, 0, len(rows))
	for i := range rows {
		result = append(result, converter.ToImMessageResult(rows[i]))
	}

	return result, nil
}
