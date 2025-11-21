package database

import (
	"project/pkg/config"
	"project/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioEntity *Minio

type Minio struct {
	//组件名称
	name string
	//地址
	endpoint string
	//账号
	accessKeyID string
	//密码
	secretAccessKey string
	//是否使用https
	source bool
	//是否初始化
	isInit bool
	//minio链接
	minioClient *minio.Client
}

func newMinio() *Minio {
	return &Minio{
		name:            "minio",
		isInit:          false,
		endpoint:        config.Get().Minio.Endpoint,
		accessKeyID:     config.Get().Minio.AccessKey,
		secretAccessKey: config.Get().Minio.SecretKey,
		source:          config.Get().Minio.Source,
	}
}

func GetMinio() *Minio {

	if minioEntity != nil {
		return minioEntity
	}

	logger.Sugar.Infof("\t\t[component] minio is initiating...")

	//判断配置文件是否加载成功
	if config.Get() == nil || config.Get().Minio == nil {
		logger.Sugar.Errorf("\t[component] minio config load failed")
		return nil
	}

	minioEntity = newMinio()

	return minioEntity
}

func (e *Minio) GetClient() *minio.Client {
	return e.minioClient
}

func (e *Minio) GetName() string {
	return e.name
}

func (e *Minio) InitComponent() bool {

	if e.isInit {
		logger.Sugar.Infof("\t[component] %s is inited", e.name)
		return true
	}

	var err error

	//初始化
	e.minioClient, err = minio.New(e.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(e.accessKeyID, e.secretAccessKey, ""),
		Secure: e.source,
	})
	if err != nil {
		logger.Sugar.Infof("\t\t[component] %s init failed: %s", e.name, err)
		return false
	}

	e.isInit = true

	logger.Sugar.Infof("\t\t[component] %s init success", e.name)
	return true
}

func (e *Minio) IsInitialize() bool {
	return e.isInit
}
