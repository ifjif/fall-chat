package filter

import "zzc/fall-chat/src/utils"

var FriendshipRequestFilterDefinition = FilterDefinition{
	Fields: []FilterField{
		{Column: "status", Op: In},
	},
}

type FriendshipRequestFilter struct {
	Status []int
}

func (f *FriendshipRequestFilter) IsNil() bool {
	return f == nil
}

func (f *FriendshipRequestFilter) ToMap() map[string]any {
	return map[string]any{
		"status": utils.ToAnySlice(f.Status),
	}
}
