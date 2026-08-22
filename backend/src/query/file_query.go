package query

type FileInstanceCheckQuery struct {
	ID       int64
	Username string
	FileID   string
}

type FileChunkStatusQuery struct {
	ID       int64
	Username string
	FileID   string
}

type FileFetchQuery struct {
	ID       int64
	Username string
	File     string
}
