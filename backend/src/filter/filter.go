package filter

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type CombineFilter struct {
	Combine []FilterBase
}

func (c *CombineFilter) ToMap() map[string]any {
	return nil
}

func (c *CombineFilter) IsNil() bool {
	return c == nil
}

type FilterBase interface {
	ToMap() map[string]any
	IsNil() bool
}

func ToMap(f FilterBase) map[string]any {
	if f == nil || f.IsNil() {
		return nil
	}

	return f.ToMap()
}

var Empty = FilterDefinition{}

type FilterOp string

const (
	Eq      FilterOp = "="
	Ne      FilterOp = "!="
	Gt      FilterOp = ">"
	Gte     FilterOp = ">="
	Lt      FilterOp = "<"
	Lte     FilterOp = "<="
	Like    FilterOp = "LIKE"
	In      FilterOp = "IN"
	Between FilterOp = "BETWEEN"
)

func BuildFilterScope(def FilterDefinition, values map[string]any, handleColumn func(f string) string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sql, args, err := def.Where(values, handleColumn)
		if args == nil || err != nil {
			return db
		}

		return db.Where(sql, args...)
	}
}

type ListFilterOptions struct {
	Filter    map[string]any
	FilterDef FilterDefinition
}

type FilterField struct {
	Column string
	Op     FilterOp
}

type FilterDefinition struct {
	Fields []FilterField
}

func (fd FilterDefinition) Where(values map[string]any, handleColumn func(f string) string) (string, []any, error) {
	if handleColumn == nil {
		handleColumn = func(f string) string {
			return f
		}
	}
	if len(values) == 0 {
		return "", nil, nil
	}

	var parts []string
	var args []any

	for _, f := range fd.Fields {
		v, ok := values[f.Column]
		if !ok || v == nil {
			continue
		}

		column := handleColumn(f.Column)

		switch f.Op {
		case In:
			slice, ok := v.([]any)
			if !ok || len(slice) == 0 {
				continue
			}

			placeholders := make([]string, 0, len(slice))
			for range slice {
				placeholders = append(placeholders, "?")
			}
			parts = append(parts, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
			args = append(args, slice...)

		case Between:
			pair, ok := v.([]any)
			if !ok || len(pair) != 2 {
				continue
			}
			parts = append(parts, fmt.Sprintf("%s BETWEEN ? AND ?", column))
			args = append(args, pair[0], pair[1])

		case Like:
			parts = append(parts, fmt.Sprintf("%s LIKE ?", column))
			args = append(args, v)

		default:
			parts = append(parts, fmt.Sprintf("%s %s ?", column, f.Op))
			args = append(args, v)
		}
	}

	if len(parts) == 0 {
		return "", nil, nil
	}

	return strings.Join(parts, " AND "), args, nil
}
