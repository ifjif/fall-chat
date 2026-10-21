package query

type ConvMembListBrief4Conv struct {
	ID       int64
	Username string
	ConvID   string
}

type ConvMembListBriefBatch struct {
	ID        int64
	Username  string
	ConvID    string
	MemberIDs []string
}
