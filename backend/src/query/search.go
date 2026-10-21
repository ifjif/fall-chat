package query

type UserOrGroup struct {
	ID         int64
	Username   string
	Name       string
	NextCursor string
}
