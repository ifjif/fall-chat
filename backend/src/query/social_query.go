package query

type SocialSearchFriendRequestQuery struct {
	ID           int64
	Username     string
	FilterStatus []int
}
type SocialSearchGroupRequestQuery struct {
	ID           int64
	Username     string
	FilterStatus []int
}

type SocialFriendshipListQuery struct {
	ID       int64
	Username string
}
