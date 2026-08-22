package response

import "github.com/minio/minio-go/v7"

type FileConfigResponse struct {
	ID        string `json:"id"`
	MaxSize   int64  `json:"max_size"`
	ChunkSize int64  `json:"chunk_size"`
}

type FileFetchResult struct {
	Object *minio.Object
	Stat   minio.ObjectInfo
}
