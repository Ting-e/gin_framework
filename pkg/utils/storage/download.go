package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	local_config "project/pkg/config"
	"project/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// DownloadFile 下载文件
//
// 下载文件到本地
//
//	fileName：文件在存储桶中的key
func DownloadFile(fileName string) ([]byte, error) {

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
		return nil, fmt.Errorf("加载Storage配置失败: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	resp, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(localCfg.Storage.BucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {

		logger.Sugar.Error(err)
		return nil, nil
	}
	defer resp.Body.Close()

	// 读取内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应体失败: %v", err)
	}

	return body, nil
}
