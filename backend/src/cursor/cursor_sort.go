package cursor

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var EmptySortDef = SortDefinition{}

type Direction string

const (
	ASC  Direction = "asc"
	DESC Direction = "desc"
)

func BuildSortOrderScope(def SortDefinition, handleColumn func(f string) string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		orderby := def.OrderBy(handleColumn)
		if orderby == "" {
			return db
		}
		return db.Order(orderby)
	}
}

func BuildSortWhereScope(def SortDefinition, values map[string]any, handleColumn func(f string) string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		sql, args, err := def.Where(values, handleColumn)
		if args == nil || err != nil {
			return db
		}

		return db.Where(sql, args...)
	}
}

func BuildSortScope(def SortDefinition, values map[string]any, handleColumn func(f string) string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = BuildSortOrderScope(def, handleColumn)(db)
		db = BuildSortWhereScope(def, values, handleColumn)(db)
		return db
	}
}

type ListSortOptions struct {
	Sort    map[string]any
	SortDef SortDefinition
}

type SortField struct {
	Field     string    `json:"field"`
	Direction Direction `json:"direction"`
}

type SortDefinition struct {
	Fields []SortField
}

// 校验排序定义
func (sd SortDefinition) Validate() error {
	if len(sd.Fields) == 0 {
		return fmt.Errorf("sort definition must have at least on field")
	}

	for _, f := range sd.Fields {
		if f.Field == "" {
			return fmt.Errorf("sort field name must not be empty")
		}

		if f.Direction != ASC && f.Direction != DESC {
			return fmt.Errorf("invalid direction for field %s: %s", f.Field, f.Direction)
		}
	}

	return nil
}

// 生成 order by 从句
func (sd SortDefinition) OrderBy(handleField func(field string) string) string {
	if handleField == nil {
		handleField = func(field string) string { return field }
	}
	parts := make([]string, 0, len(sd.Fields))
	for _, f := range sd.Fields {
		fstr := handleField(f.Field)
		parts = append(parts, fmt.Sprintf("%s %s", fstr, strings.ToUpper(string(f.Direction))))
	}
	return strings.Join(parts, ", ")
}

// 生成 Where从句
// 只加入values中存在的
//
//	  hot < xx
//		OR (hot = xx and create_time < xx)
//		OR (hot = xx and create_time = xx and id < xx)
func (sd SortDefinition) Where(values map[string]any, handleField func(field string) string) (string, []any, error) {
	if err := sd.Validate(); err != nil {
		return "", nil, err
	}

	if len(values) == 0 {
		return "", nil, nil
	}

	if handleField == nil {
		handleField = func(field string) string { return field }
	}

	// 过滤出 values中存在的fields
	var fields []SortField
	for i := range sd.Fields {
		name := sd.Fields[i].Field
		if _, ok := values[name]; ok {
			fields = append(fields, sd.Fields[i])
		}
	}

	var parts []string
	var args []any

	for i, field := range fields {
		var eq []string
		for j := range i {
			prev := sd.Fields[j]
			fstr := handleField(prev.Field)
			eq = append(eq, fmt.Sprintf("%s = ?", fstr))
			args = append(args, values[prev.Field])
		}

		op := ">"
		if field.Direction == DESC {
			op = "<"
		}

		ftr := handleField(field.Field)
		cond := fmt.Sprintf("%s %s ?", ftr, op)
		args = append(args, values[field.Field])

		if len(eq) > 0 {
			parts = append(parts, fmt.Sprintf("(%s AND %s)", strings.Join(eq, "AND"), cond))
		} else {
			parts = append(parts, cond)
		}
	}

	if len(parts) == 0 {
		return "", nil, nil
	}

	return strings.Join(parts, " OR "), args, nil
}
