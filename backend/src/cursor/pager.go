package cursor

import (
	"context"
)

// 游标分页器
type CursorPager[C CursorBase, T any, R any] struct {
	Limit       int
	DecodeNext  func(s string) (C, error)
	GetSortDef  func(key string) (SortDefinition, error)
	Query       func(ctx context.Context, sortValue map[string]any, sortDef SortDefinition, limit int) ([]T, error)
	BuildCursor func(next C, rows []T, lastIdx int, hasNext bool) C
	ToResult    func(T) R
}

func (p *CursorPager[C, T, R]) List(
	ctx context.Context,
	nextStr string,
	sortKey string,
) ([]R, *string, bool, error) {
	limit := p.Limit
	if limit <= 0 {
		limit = 20
	}

	next, err := p.DecodeNext(nextStr)
	if err != nil {
		return nil, nil, false, err
	}

	if NoMore(next) {
		return nil, &nextStr, false, nil
	}

	sortDef, err := p.GetSortDef(sortKey)
	if err != nil {
		return nil, nil, false, err
	}

	sortValue := ToMap(next)

	rows, err := p.Query(ctx, sortValue, sortDef, limit+1)
	if err != nil {
		return nil, nil, false, err
	}

	newCursor := Update(next, len(rows), limit, func(lastIdx int, hasNext bool) C {
		return p.BuildCursor(next, rows, lastIdx, hasNext)
	})

	hasNext := HashNext(newCursor)
	if hasNext {
		rows = rows[:limit]
	}

	result := make([]R, 0, len(rows))
	for i := range rows {
		result = append(result, p.ToResult(rows[i]))
	}

	encoded, err := Encode(newCursor)
	if err != nil {
		return nil, nil, false, err
	}

	return result, &encoded, hasNext, nil
}
