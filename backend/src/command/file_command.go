package command

import "mime/multipart"

type FileUploadChunkCommand struct {
	ID       int64
	Username string
	FileID   string
	Index    int
	Chunk    *multipart.FileHeader
}

type FileUploadMergeCommand struct {
	ID         int64
	Username   string
	Filename   string
	FileID     string
	TotalChunk int
	Chunks     []int
}

type FileUploadCommand struct {
	ID       int64
	Username string
	FileID   string
	File     *multipart.FileHeader
}
