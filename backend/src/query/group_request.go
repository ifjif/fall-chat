package query

type ListGroupRequest struct {
	ID           int64
	Username     string
	FilterStatus []int
}
