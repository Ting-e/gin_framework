package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project/pkg/utils/snowflake"

	"project/internal/model"

	"os"
	local_config "project/pkg/config"
	"project/pkg/logger"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var defaultService *Service

type Service struct {
	snowflake *snowflake.Worker
}

func GetService() *Service {
	return defaultService
}

func init() {
	defaultService = &Service{
		snowflake: snowflake.NewWorker(snowflake.WorkerID, snowflake.WataCenterID),
	}
}

// 连接redis
func ConnectRedis(db int) (*redis.Client, error) {

	cfg := local_config.Get()

	rdb := redis.NewClient(&redis.Options{
		// Addr: "localhost:63794",
		Addr:     cfg.Db.Redis.Addr,
		Password: cfg.Db.Redis.Password,
		DB:       db,
	})

	// 检查客户端是否可以Ping Redis服务器以验证连接
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("无法连接到Redis服务器: %v", err)
	}

	return rdb, nil
}

// 生成随机ID（UUID）
func GenerateID(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("长度必须大于0")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("未能生成 UUID: %w", err)
	}

	// 使用 base64 编码缩短 ID 长度
	b := make([]byte, base64.URLEncoding.EncodedLen(len(id)))
	base64.URLEncoding.Encode(b, id[:])

	// 截取前 length 个字符
	shortenedID := string(b)[:length]

	return shortenedID, nil
}

// 生成随机ID（小写+数字）
func GenerateKey(length int) string {
	letters := "abcdefghijkmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	for i, v := range b {
		b[i] = letters[v%byte(len(letters))]
	}
	return string(b)
}

// 发送post请求
func PostRequest(url string, requestBody interface{}) (*model.Response, error) {
	// 将请求体转换为JSON格式
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("无法封送请求正文: %w", err)
	}
	// 创建HTTP请求
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("无法创建HTTP请求: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法执行HTTP请求: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("收到未成功的状态码: %d, 响应体: %s", resp.StatusCode, bodyBytes)
	}
	// 解析响应体为Response结构体
	res := new(model.Response)
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("无法解码响应正文: %w", err)
	}
	return res, nil
}

// 上传本地文件并删除（后端直传）
func UplodaLocalFile(filePath string) error {

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// 构建上传文件名
	s := strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond)))
	Key := "testFloder/" + s + ".xlsx"

	ctx := context.Background()
	local_cfg := local_config.Get()

	// 初始化上传凭证
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(local_cfg.Minio.AccessKey, local_cfg.Minio.SecretKey, "")),
		config.WithRegion(local_cfg.Minio.Region),
		config.WithBaseEndpoint(local_cfg.Minio.Endpoint),
	)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// 上传
	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(local_cfg.Minio.BucketName),
		Key:    aws.String(Key),
		Body:   file,
	})
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}
	file.Close()

	// 删除本地文件
	err = os.Remove(filePath)
	if err != nil {
		logger.Sugar.Error(err)
		return err
	}

	return nil
}

// 获取上传预签名（预签名上传）
func PresignPutObjectUrl(filePath string) (string, error) {

	// 构建上传文件名
	s := strconv.Itoa(int(time.Now().UnixNano() / int64(time.Millisecond)))
	key := "testFloder/" + s + ".xlsx"

	ctx := context.Background()
	local_cfg := local_config.Get()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(local_cfg.Minio.AccessKey, local_cfg.Minio.SecretKey, "")),
		config.WithRegion(local_cfg.Minio.Region),
		config.WithBaseEndpoint("https://"+local_cfg.Minio.Endpoint),
	)
	if err != nil {
		logger.Sugar.Error(err)
		return "", err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(client)

	response, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(local_cfg.Minio.BucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(1000 * int64(time.Second))
	})
	if err != nil {
		logger.Sugar.Error(err)
		return "", err
	}

	presignedURL := response.URL

	return presignedURL, nil
}
