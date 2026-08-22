package filestore

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

const (
	ChunksFolder = "chunks"
	FilesFolder  = "files"
	ChunkSize    = 5 * 1024 * 1024
)

var FileStore *Store

type CommonOptions struct {
	Bucket string
	Ctx    context.Context
}

type StoreOptions struct {
	CommonOptions
	GetObjectOptions   *minio.GetObjectOptions
	PutObjectOptions   *minio.PutObjectOptions
	ListObjectsOptions *minio.ListObjectsOptions
}

type Store struct {
	minioClient  *minio.Client
	bucket       string
	chunksFolder string
	filesFolder  string
}

func (s *Store) Bucket() string {
	return s.bucket
}

func (s *Store) ChunksFolder() string {
	return s.chunksFolder
}

func (s *Store) FilesFolder() string {
	return s.filesFolder
}

func (s *Store) CommonOptions(options StoreOptions) *CommonOptions {
	comm := options.CommonOptions
	if comm.Bucket == "" {
		comm.Bucket = s.bucket
	}
	if comm.Ctx == nil {
		comm.Ctx = context.Background()
	}

	return &comm
}

func (s *Store) GetObject(objectName string, options StoreOptions) (*minio.Object, error) {
	comm := s.CommonOptions(options)
	getObjectOptions := minio.GetObjectOptions{}
	if options.GetObjectOptions != nil {
		getObjectOptions = *options.GetObjectOptions
	}
	return s.minioClient.GetObject(comm.Ctx, comm.Bucket, objectName, getObjectOptions)
}

func (s *Store) PutObject(objectName string, size int64, reader io.Reader, options StoreOptions) (minio.UploadInfo, error) {
	comm := s.CommonOptions(options)
	putObjectOptions := minio.PutObjectOptions{}
	if options.PutObjectOptions != nil {
		putObjectOptions = *options.PutObjectOptions
	}

	info, err := s.minioClient.PutObject(comm.Ctx, comm.Bucket, objectName, reader, size, putObjectOptions)
	return info, err
}

func (s *Store) ComposeObject(dst minio.CopyDestOptions, srcs []minio.CopySrcOptions, options StoreOptions) (minio.UploadInfo, error) {
	comm := s.CommonOptions(options)

	info, err := s.minioClient.ComposeObject(comm.Ctx, dst, srcs...)

	return info, err
}

func (s *Store) ListObjects(options StoreOptions) <-chan minio.ObjectInfo {
	comm := s.CommonOptions(options)

	listObjectsOptions := minio.ListObjectsOptions{}
	if options.ListObjectsOptions != nil {
		listObjectsOptions = *options.ListObjectsOptions
	}

	return s.minioClient.ListObjects(comm.Ctx, comm.Bucket, listObjectsOptions)
}
