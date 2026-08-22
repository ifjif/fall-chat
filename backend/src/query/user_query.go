package query

type SearchQuery struct {
	ID         int64
	Username   string
	Name       string
	NextCursor string
}
