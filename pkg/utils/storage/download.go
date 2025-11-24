package storage

import (
	"context"
	"fmt"

	"os"
	local_config "project/pkg/config"
	"project/pkg/logger"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DownloadFile 下载文件
func DownloadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 构建上传文件名
	timestamp := strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond)))
	key := "testFolder/" + timestamp + ".xlsx"

	ctx := context.Background()
	localCfg := local_config.Get()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			localCfg.Storage.AccessKey,
			localCfg.Storage.SecretKey,
			"",
		)),
		config.WithRegion(localCfg.Storage.Region),
		config.WithBaseEndpoint(localCfg.Storage.Endpoint),
	)
	if err != nil {
		logger.Sugar.Error("加载Storage配置失败:", err)
		return fmt.Errorf("加载Storage配置失败: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(localCfg.Storage.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		logger.Sugar.Error("上传文件到存储桶失败:", err)
		return fmt.Errorf("上传文件到存储桶失败: %w", err)
	}

	// 删除本地文件
	err = os.Remove(filePath)
	if err != nil {
		logger.Sugar.Error("删除本地文件失败:", err)
		return fmt.Errorf("删除本地文件失败: %w", err)
	}

	return nil
}
