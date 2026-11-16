package command

import "mime/multipart"

type FileUploadChunk struct {
	UserID   int64
	Username string
	FileID   string
	Index    int
	Chunk    *multipart.FileHeader
}

type FileUploadMerge struct {
	UserID     int64
	Username   string
	Filename   string
	FileID     string
	TotalChunk int
	Chunks     []int
}

type FileUpload struct {
	UserID   int64
	Username string
	FileID   string
	File     *multipart.FileHeader
}
