package request

import (
	"mime/multipart"
)

type FileFetch struct {
	File string `query:"file" validate:"required"`
}

type FileInstanceCheck struct {
	FileID string `json:"file_id" validate:"required"`
}
type FileChunkStatus struct {
	FileID string `json:"file_id" validate:"required"`
}

type FileUploadChunk struct {
	FileID string                `form:"file_id" validate:"required"`
	Index  int                   `form:"index" validate:"min=0"`
	Chunk  *multipart.FileHeader `form:"chunk" validate:"required"`
}

type FileUpload struct {
	FileID string                `form:"file_id" validate:"required"`
	File   *multipart.FileHeader `form:"file" validate:"required"`
}

type FileUploadMerge struct {
	FileID     string `json:"file_id" validate:"required"`
	Filename   string `json:"filename" validate:"required"`
	TotalChunk int    `json:"total_chunk" validate:"required,min=0"`
	Chunks     []int  `json:"chunks" validate:"required"`
}
