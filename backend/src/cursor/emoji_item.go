package cursor

var EmojiItemSortDefinitions = map[string]SortDefinition{
	"sort_no_asc": {
		Fields: []SortField{
			{Field: "sort_no", Direction: ASC},
			{Field: "id", Direction: ASC},
		},
	},

	"sort_no_desc": {
		Fields: []SortField{
			{Field: "sort_no", Direction: DESC},
			{Field: "id", Direction: DESC},
		},
	},
}

type EmojiItemCursor struct {
	Base
	SortNo int    `json:"sort_no"`
	ID     uint64 `json:"id"`
}

func (c *EmojiItemCursor) ToMap() map[string]any {
	return map[string]any{
		"sort_no": c.SortNo,
		"id":      c.ID,
	}
}

func (c *EmojiItemCursor) IsNil() bool {
	return c == nil
}
