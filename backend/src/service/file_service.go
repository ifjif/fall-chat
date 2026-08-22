package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/filestore"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

const (
	DEFAULT_MAX_SIZE   = 100 * 1024 * 1024 // 100M
	DEFAULT_CHUNK_SIZE = 20 * 1024 * 1024  // 20M
)

type FileService struct {
	BaseService
	store *filestore.Store
}

func NewFileService(name string) *FileService {
	bs := BaseService{name: name}
	return &FileService{
		BaseService: bs,
	}
}

func (fs *FileService) Inject() {
	fs.BaseService.Inject()
	fs.store = filestore.FileStore
}

func (fs *FileService) add(maxSize int64, chunkSize int64) (*model.FileConfig, error) {
	id, err := fs.sf.NextID()
	if err != nil {
		fs.LogWarn("add", err, "reason", "添加文件配置信息失败")
	}
	publicID := fs.sf.Encode(id)

	fc := &model.FileConfig{
		MaxSize:   maxSize,
		ChunkSize: chunkSize,
	}
	fc.ID = uint64(id)
	fc.PublicID = publicID

	result := fs.db.Create(fc)
	if result.Error != nil {
		fs.LogWarn("add", err, "reason", "创建一条记录失败")
		return nil, result.Error
	}

	return fc, nil
}

func (fs *FileService) Config() (*model.FileConfig, error) {
	var fc model.FileConfig
	result := fs.db.Take(&fc)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fc, err := fs.add(DEFAULT_MAX_SIZE, DEFAULT_CHUNK_SIZE)
			if err != nil {
				return nil, err
			}
			return fc, nil
		}

		fs.LogWarn("Config", result.Error, "reason", "获取文件配置信息失败")
		return nil, result.Error
	}

	return &fc, nil
}

func (fs *FileService) Fetch(qu *query.FileFetchQuery) (*response.FileFetchResult, error) {
	object, err := fs.store.GetObject(qu.File, filestore.StoreOptions{})
	if err != nil {
		fs.LogError("Fetch", err,
			"username", qu.Username,
			"reason", "minio 获取文件失败",
		)

		return nil, err
	}

	stat, err := object.Stat()
	if err != nil {
		fs.LogError("Fetch", err, "username", qu.Username,
			"reason", "获取文件stat失败",
		)

		return nil, err
	}

	result := &response.FileFetchResult{
		Object: object,
		Stat:   stat,
	}

	return result, err
}

func (fs *FileService) Upload(cmd *command.FileUploadCommand) (string, error) {
	targetFile := cmd.File
	filename := fmt.Sprintf("%s/%s/%s", fs.store.FilesFolder(), cmd.FileID, targetFile.Filename)

	file, err := targetFile.Open()
	if err != nil {
		fs.LogWarn("Upload", err,
			"username", cmd.Username,
			"reason", "文件打开失败",
		)
	}
	defer file.Close()

	_, err = fs.store.PutObject(filename, targetFile.Size, file, filestore.StoreOptions{
		PutObjectOptions: &minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	})
	if err != nil {
		fs.LogError("Upload", err,
			"username", cmd.Username,
			"reason", "minio 上传文件失败",
		)
		return "", err
	}

	return filename, nil
}

func (fs *FileService) UploadChunk(cmd *command.FileUploadChunkCommand) error {
	chunkName := fmt.Sprintf("%s/%s/%d", fs.store.ChunksFolder(), cmd.FileID, cmd.Index)

	file, err := cmd.Chunk.Open()
	if err != nil {
		fs.LogWarn("UploadChunk", err,
			"username", cmd.Username,
			"reason", "打开文件失败",
		)

		return err
	}
	defer file.Close()

	_, err = fs.store.PutObject(chunkName, cmd.Chunk.Size, file, filestore.StoreOptions{
		PutObjectOptions: &minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	})
	if err != nil {
		fs.LogError("UploadChunk", err, "reason", "上传到minio 错误")
		return err
	}
	return nil
}

func (fs *FileService) UploadMerge(cmd *command.FileUploadMergeCommand) (string, error) {
	if cmd.TotalChunk != len(cmd.Chunks) {
		fs.LogWarn("UploadMerge", errs.FileErrChunkCountInvalid, "username", cmd.Username)
		return "", errs.FileErrChunkCountInvalid
	}

	var srcs []minio.CopySrcOptions
	sort.Ints(cmd.Chunks)
	for i := range cmd.Chunks {
		name := fmt.Sprintf("%s/%s/%d", fs.store.ChunksFolder(), cmd.FileID, cmd.Chunks[i])
		src := minio.CopySrcOptions{
			Bucket: fs.store.Bucket(),
			Object: name,
		}
		srcs = append(srcs, src)
	}

	filename := fmt.Sprintf("%s/%s/%s", fs.store.FilesFolder(), cmd.FileID, cmd.Filename)
	dist := minio.CopyDestOptions{
		Bucket: fs.store.Bucket(),
		Object: filename,
	}

	_, err := fs.store.ComposeObject(dist, srcs, filestore.StoreOptions{})
	if err != nil {
		fs.LogError("UploadMerge", err,
			"username", cmd.Username,
			"reason", "minio compose object 失败",
		)
		return "", err
	}
	return filename, nil
}

func (fs *FileService) ChunkStatus(qu *query.FileChunkStatusQuery) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objectCh := fs.store.ListObjects(filestore.StoreOptions{
		CommonOptions: filestore.CommonOptions{
			Ctx: ctx,
		},
		ListObjectsOptions: &minio.ListObjectsOptions{
			Prefix:    fs.store.ChunksFolder() + "/" + qu.FileID + "/",
			Recursive: true,
		},
	})

	keys := []string{}

	for info := range objectCh {
		if info.Err != nil {
			fs.LogError("ChunksFolder", info.Err,
				"username", qu.Username,
				"reason", "minio获取已上传chunk失败")

			return nil, info.Err
		}
		key := info.Key
		names := strings.Split(key, "/")
		keys = append(keys, names[len(names)-1])
	}

	return keys, nil
}

func (fs *FileService) InstanceCheck(qu *query.FileInstanceCheckQuery) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objecCh := fs.store.ListObjects(filestore.StoreOptions{
		CommonOptions: filestore.CommonOptions{
			Ctx: ctx,
		},
		ListObjectsOptions: &minio.ListObjectsOptions{
			Prefix:    fs.store.FilesFolder() + "/" + qu.FileID + "/",
			Recursive: true,
		},
	})

	select {
	case obj, ok := <-objecCh:
		if !ok {
			return "", nil
		}
		if obj.Err != nil {
			fs.LogError("InstanceCheck", obj.Err,
				"usernae", qu.Username,
				"reason", "从minio list object失败",
			)

			return "", obj.Err
		}

		return obj.Key, nil
	}
}
