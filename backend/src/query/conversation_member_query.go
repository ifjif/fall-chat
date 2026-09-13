package query

type ConvMembListBrief4ConvQuery struct {
	ID       int64
	Username string
	ConvID   string
}

type ConvMembListBriefBatchQuery struct {
	ID        int64
	Username  string
	ConvID    string
	MemberIDs []string
}
