package filestore

import (
	"log"
	"strconv"
	"zzc/fall-chat/src/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	Bucketname      = ""
	EndPoint        = ""
	AccessKeyID     = ""
	SecretAccessKey = ""
	UseSSL          = false
)

func init() {
	Bucketname, _ := utils.GetEnv("MINIO_BUCKETNAME")
	EndPoint, _ := utils.GetEnv("MINIO_ENDPOINT")
	AccessKeyID, _ := utils.GetEnv("MINIO_ACCESSKEYID")
	SecretAccessKey, _ := utils.GetEnv("MINIO_SECRETACCESSKEY")
	UseSSLStr, _ := utils.GetEnv("MINIO_USESSL")
	UseSSL, err := strconv.ParseBool(UseSSLStr)
	if err != nil {
		utils.LogError("minio-client UseSSl解析错误，采用默认禁用")
	}

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
