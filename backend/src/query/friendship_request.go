package query

type ListFriendRequest struct {
	ID           int64
	Username     string
	FilterStatus []int
}
