package query

import "zzc/fall-chat/src/filter"

type ListGroupRequest struct {
	UserID   int64
	Username string
	Next     string
	SortKey  string
	Filter   *filter.GroupRequestFilter
}
