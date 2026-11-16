package query

import "zzc/fall-chat/src/filter"

type ListFriendRequest struct {
	UserID   uint64
	Username string
	Filter   *filter.FriendshipRequestFilter
	Next     string
	SortKey  string
}
