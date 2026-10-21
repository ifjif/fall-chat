package command

import "mime/multipart"

type FileUploadChunk struct {
	ID       int64
	Username string
	FileID   string
	Index    int
	Chunk    *multipart.FileHeader
}

type FileUploadMerge struct {
	ID         int64
	Username   string
	Filename   string
	FileID     string
	TotalChunk int
	Chunks     []int
}

type FileUpload struct {
	ID       int64
	Username string
	FileID   string
	File     *multipart.FileHeader
}
