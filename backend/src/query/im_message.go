package query

type ListIMMessage struct {
	UserID   int64
	Username string
	ConvID   int64
	Next     string
	SortKey  string
}

type PatchListIMMessage struct {
	UserID   uint64
	Username string
	ConvID   int64
	MinSeq   int64
	MaxSeq   int64
}
