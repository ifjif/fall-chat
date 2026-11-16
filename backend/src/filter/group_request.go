package filter

import "zzc/fall-chat/src/utils"

var GroupRequestFilterDefinition = FilterDefinition{
	Fields: []FilterField{
		{Column: "status", Op: In},
	},
}

type GroupRequestFilter struct {
	Status []int
}

func (f *GroupRequestFilter) IsNil() bool {
	return f == nil
}

func (f *GroupRequestFilter) ToMap() map[string]any {
	return map[string]any{
		"status": utils.ToAnySlice(f.Status),
	}
}
