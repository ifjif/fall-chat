package model

type FileConfig struct {
	Model
	MaxSize   int64
	ChunkSize int64
}
