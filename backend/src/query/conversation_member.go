package query

type ConvMembListBrief4Conv struct {
	UserID   int64
	Username string
	ConvID   int64
}

type ConvMembListBriefBatch struct {
	UserID     int64
	Username   string
	ConvID     int64
	MemberIDs  []int64
	MemberPIDs []string
}
