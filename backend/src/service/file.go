package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
	"zzc/fall-chat/src/command"
	"zzc/fall-chat/src/data"
	"zzc/fall-chat/src/data/store"
	"zzc/fall-chat/src/dto/response"
	"zzc/fall-chat/src/errs"
	"zzc/fall-chat/src/filestore"
	"zzc/fall-chat/src/model"
	"zzc/fall-chat/src/query"

	"github.com/minio/minio-go/v7"
)

const (
	DEFAULT_MAX_SIZE   = 100 * 1024 * 1024 // 100M
	DEFAULT_CHUNK_SIZE = 20 * 1024 * 1024  // 20M
)

type File struct {
	BaseService
	fstore *filestore.Store
	store  store.FileStore
}

func NewFile(name string) *File {
	bs := BaseService{name: name}
	return &File{
		BaseService: bs,
	}
}

func (fs *File) Inject() {
	fs.BaseService.Inject()
	fs.fstore = filestore.FileStore
	fs.store = data.GetBaseDao()
}

func (fs *File) add(ctx context.Context, maxSize int64, chunkSize int64) (*model.FileConfig, error) {
	id, err := fs.sf.NextID()
	if err != nil {
		return nil, fs.Wrap("生成文件配置信息ID", err)
	}
	publicID := fs.sf.Encode(id)

	fc := &model.FileConfig{
		MaxSize:   maxSize,
		ChunkSize: chunkSize,
	}
	fc.ID = uint64(id)
	fc.PublicID = publicID

	err = fs.store.FileRepos().File.CreateOne(ctx, *fc)
	if err != nil {
		return nil, fs.Wrap("创建配置文件记录", err)
	}

	return fc, nil
}

func (fs *File) Config(ctx context.Context) (*model.FileConfig, error) {
	fc, err := fs.store.FileRepos().File.FindOne(ctx)
	if err != nil {
		return nil, fs.Wrap("获取文件配置信息", err)
	}

	if fc != nil {
		return fc, nil
	}

	fc, err = fs.add(ctx, DEFAULT_MAX_SIZE, DEFAULT_CHUNK_SIZE)
	if err != nil {
		return nil, fs.Wrap("添加文件配置信息", err)
	}
	return fc, nil
}

func (fs *File) Fetch(qu *query.FileFetch) (*response.FileFetchResult, error) {
	object, err := fs.fstore.GetObject(qu.File, filestore.StoreOptions{})
	if err != nil {
		return nil, fs.Wrap("minio 获取文件", err, "username", qu.Username)
	}

	stat, err := object.Stat()
	if err != nil {
		return nil, fs.Wrap("获取文件stat", err, "username", qu.Username)
	}

	result := &response.FileFetchResult{
		Object: object,
		Stat:   stat,
	}

	return result, err
}

func (fs *File) Upload(cmd *command.FileUpload) (string, error) {
	targetFile := cmd.File
	filename := fmt.Sprintf("%s/%s/%s", fs.fstore.FilesFolder(), cmd.FileID, targetFile.Filename)

	file, err := targetFile.Open()
	if err != nil {
		return "", fs.Wrap("文件打开", err, "username", cmd.Username)
	}
	defer file.Close()

	_, err = fs.fstore.PutObject(filename, targetFile.Size, file, filestore.StoreOptions{
		PutObjectOptions: &minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	})
	if err != nil {
		return "", fs.Wrap("o 上传文件", err, "username", cmd.Username)
	}

	return filename, nil
}

func (fs *File) UploadChunk(cmd *command.FileUploadChunk) error {
	chunkName := fmt.Sprintf("%s/%s/%d", fs.fstore.ChunksFolder(), cmd.FileID, cmd.Index)

	file, err := cmd.Chunk.Open()
	if err != nil {
		return fs.Wrap("打开文件", err,
			"username", cmd.Username)
	}
	defer file.Close()

	_, err = fs.fstore.PutObject(chunkName, cmd.Chunk.Size, file, filestore.StoreOptions{
		PutObjectOptions: &minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	})
	if err != nil {
		return fs.Wrap("上传到minio ", err)
	}
	return nil
}

func (fs *File) UploadMerge(cmd *command.FileUploadMerge) (string, error) {
	if cmd.TotalChunk != len(cmd.Chunks) {
		return "", fs.Wrap("文件chunk数量比对", errs.FileErrChunkCountInvalid, "username", cmd.Username)
	}

	var srcs []minio.CopySrcOptions
	sort.Ints(cmd.Chunks)
	for i := range cmd.Chunks {
		name := fmt.Sprintf("%s/%s/%d", fs.fstore.ChunksFolder(), cmd.FileID, cmd.Chunks[i])
		src := minio.CopySrcOptions{
			Bucket: fs.fstore.Bucket(),
			Object: name,
		}
		srcs = append(srcs, src)
	}

	filename := fmt.Sprintf("%s/%s/%s", fs.fstore.FilesFolder(), cmd.FileID, cmd.Filename)
	dist := minio.CopyDestOptions{
		Bucket: fs.fstore.Bucket(),
		Object: filename,
	}

	_, err := fs.fstore.ComposeObject(dist, srcs, filestore.StoreOptions{})
	if err != nil {
		return "", fs.Wrap("minio compose object", err, "username", cmd.Username)
	}
	return filename, nil
}

func (fs *File) ChunkStatus(qu *query.FileChunkStatus) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objectCh := fs.fstore.ListObjects(filestore.StoreOptions{
		CommonOptions: filestore.CommonOptions{
			Ctx: ctx,
		},
		ListObjectsOptions: &minio.ListObjectsOptions{
			Prefix:    fs.fstore.ChunksFolder() + "/" + qu.FileID + "/",
			Recursive: true,
		},
	})

	keys := []string{}

	for info := range objectCh {
		if info.Err != nil {
			return nil, fs.Wrap("minio获取已上传chunk", info.Err, "username", qu.Username)
		}
		key := info.Key
		names := strings.Split(key, "/")
		keys = append(keys, names[len(names)-1])
	}

	return keys, nil
}

func (fs *File) InstanceCheck(qu *query.FileInstanceCheck) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objecCh := fs.fstore.ListObjects(filestore.StoreOptions{
		CommonOptions: filestore.CommonOptions{
			Ctx: ctx,
		},
		ListObjectsOptions: &minio.ListObjectsOptions{
			Prefix:    fs.fstore.FilesFolder() + "/" + qu.FileID + "/",
			Recursive: true,
		},
	})

	select {
	case obj, ok := <-objecCh:
		if !ok {
			return "", nil
		}
		if obj.Err != nil {
			return "", fs.Wrap("从minio list object", obj.Err, "usernae", qu.Username)
		}

		return obj.Key, nil
	}
}
