package query

type FileInstanceCheck struct {
	UserID   int64
	Username string
	FileID   string
}

type FileChunkStatus struct {
	UserID   int64
	Username string
	FileID   string
}

type FileFetch struct {
	ID       int64
	Username string
	File     string
}
