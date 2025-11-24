package storage

import (
	"context"
	"fmt"
	"strings"

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

// UploadLocalFile 上传本地文件（后端直传）
//
// 上传文件到存储桶
//
//	filePath：文件路径（folder/file.txt）
//	uploadPath：上传路径（指定上传到存储桶中的目标路径：testFolder/）,默认为根目录
//	uploadName：上传名称（指定上传到存储桶中的目标名称：testfile.txt），默认为是时间戳
//	save：上传成功后本地文件是否保存，（true：保存；false：不保存）默认为true
func UploadLocalFile(filePath string, save bool, uploadPath string, uploadName string) error {

	if filePath == "" {
		return fmt.Errorf("filePath 不可为空")
	}

	suffix := strings.Split(filePath, ".")[len(strings.Split(filePath, "."))-1:][0]

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	if uploadName == "" {
		// 构建默认上传文件名
		timestamp := strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond)))
		uploadName = timestamp + "." + suffix
	}

	key := uploadPath + uploadName

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
		logger.Sugar.Error("上传文件到Storage失败:", err)
		return fmt.Errorf("上传文件到Storage失败: %w", err)
	}

	if !save {
		// 删除本地文件
		err = os.Remove(filePath)
		if err != nil {
			logger.Sugar.Error("删除本地文件失败:", err)
			return fmt.Errorf("删除本地文件失败: %w", err)
		}
	}

	return nil
}

// PresignPutObjectUrl 获取上传预签名（预签名上传）
//
// 上传文件到存储桶
//
//	filePath：文件路径（folder/file.txt）
//	uploadPath：上传路径（指定上传到存储桶中的目标路径：testFolder/）,默认为根目录
//	uploadName：上传名称（指定上传到存储桶中的目标名称：testfile.txt）不可为空
//	save：上传成功后本地文件是否保存，（true：保存；false：不保存）默认为true
func PresignPutObjectUrl(uploadPath string, uploadName string) (string, error) {

	if uploadName == "" {
		return "", fmt.Errorf("uploadName 不可为空")
	}

	key := uploadPath + uploadName

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
		return "", fmt.Errorf("加载Storage配置失败: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(client)

	response, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(localCfg.Storage.BucketName),
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
