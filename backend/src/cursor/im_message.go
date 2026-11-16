package cursor

var IMMessageSortDefinitions = map[string]SortDefinition{
	"seq_asc": {
		Fields: []SortField{
			{Field: "seq", Direction: ASC},
			{Field: "id", Direction: ASC},
		},
	},
	"seq_desc": {
		Fields: []SortField{
			{Field: "seq", Direction: DESC},
			{Field: "id", Direction: DESC},
		},
	},
}

type IMMessageCursor struct {
	Base
	Seq int64
	ID  uint64
}

func (c *IMMessageCursor) IsNil() bool {
	return c == nil
}

func (c *IMMessageCursor) ToMap() map[string]any {
	return map[string]any{
		"id":  c.ID,
		"seq": c.Seq,
	}
}
