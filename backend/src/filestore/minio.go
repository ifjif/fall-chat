package filestore

import (
	"log"
	"zzc/fall-chat/src/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	Bucketname      = "testbucket"
	EndPoint        = "minio:9000"
	AccessKeyID     = "Z9MEZ4g7wBsY3E2eyvJF"
	SecretAccessKey = "MHm7D7TQjAFgBssgYhHzOQleQRIttxO0cluaQGhC"
	UseSSL          = false
)

func init() {
	client, err := minio.New(EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AccessKeyID, SecretAccessKey, ""),
		Secure: UseSSL,
	})
	if err != nil {
		utils.LogError("minio-client初始失败，无法使用存储服务")
		log.Fatalln(err)
	}
	// MinioClient.TraceOn(os.Stdout)
	FileStore = &Store{
		minioClient:  client,
		bucket:       Bucketname,
		chunksFolder: ChunksFolder,
		filesFolder:  FilesFolder,
	}

	utils.LogInfo("%#v", FileStore)
}
