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

// UploadLocalFile 上传本地文件并删除（后端直传）
func UploadLocalFile(filePath string) error {
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
			localCfg.Minio.AccessKey,
			localCfg.Minio.SecretKey,
			"",
		)),
		config.WithRegion(localCfg.Minio.Region),
		config.WithBaseEndpoint(localCfg.Minio.Endpoint),
	)
	if err != nil {
		logger.Sugar.Error("加载MinIO配置失败:", err)
		return fmt.Errorf("加载MinIO配置失败: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(localCfg.Minio.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		logger.Sugar.Error("上传文件到MinIO失败:", err)
		return fmt.Errorf("上传文件到MinIO失败: %w", err)
	}

	// 删除本地文件
	err = os.Remove(filePath)
	if err != nil {
		logger.Sugar.Error("删除本地文件失败:", err)
		return fmt.Errorf("删除本地文件失败: %w", err)
	}

	return nil
}

// PresignPutObjectUrl 获取上传预签名（预签名上传）
func PresignPutObjectUrl(filename string) (string, error) {
	if filename == "" {
		filename = "temp"
	}

	timestamp := strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond)))
	key := "testFolder/" + timestamp + "_" + filename

	ctx := context.Background()
	localCfg := local_config.Get()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			localCfg.Minio.AccessKey,
			localCfg.Minio.SecretKey,
			"",
		)),
		config.WithRegion(localCfg.Minio.Region),
		config.WithBaseEndpoint(localCfg.Minio.Endpoint),
	)
	if err != nil {
		logger.Sugar.Error("加载MinIO配置失败:", err)
		return "", fmt.Errorf("加载MinIO配置失败: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(client)

	response, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(localCfg.Minio.BucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(1000 * int64(time.Second))
	})
	if err != nil {
		logger.Sugar.Error("生成预签名URL失败:", err)
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}

	return response.URL, nil
}
