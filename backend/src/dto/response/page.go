package response

var emptyPageCursor = PageCursor{
	Data:    nil,
	Next:    "",
	HasNext: false,
}

type PageCursor struct {
	Data    []any `json:"data"`
	Next    any   `json:"next"`
	HasNext bool  `json:"has_next"`
}

func EmptyPageCursor(next string) *PageCursor {
	e := emptyPageCursor
	e.Next = next
	return &e
}
