package cursor

var FriendshipRequestSortDefinitions = map[string]SortDefinition{
	"created_asc": {
		Fields: []SortField{
			{Field: "created_at", Direction: ASC},
			{Field: "id", Direction: ASC},
		},
	},
	"created_desc": {
		Fields: []SortField{
			{Field: "created_at", Direction: DESC},
			{Field: "id", Direction: DESC},
		},
	},
}

type FriendshipRequestCursor struct {
	Base
	ID uint64
}

func (c *FriendshipRequestCursor) ToMap() map[string]any {
	return map[string]any{
		"id": c.ID,
	}
}

func (c *FriendshipRequestCursor) IsNil() bool {
	return c == nil
}
