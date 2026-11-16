package cursor

var EmojiPackSortDefinitions = map[string]SortDefinition{
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

type EmojiPackCurosr struct {
	Base
	ID uint64 `json:"id"`
}

func (c *EmojiPackCurosr) ToMap() map[string]any {
	return map[string]any{
		"id": c.ID,
	}
}

func (c *EmojiPackCurosr) IsNil() bool {
	return c == nil
}
