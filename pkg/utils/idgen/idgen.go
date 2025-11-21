package idgen

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

// GenerateID 生成随机ID（UUID）
func GenerateID(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("长度必须大于0")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("未能生成 UUID: %w", err)
	}

	b := make([]byte, base64.URLEncoding.EncodedLen(len(id)))
	base64.URLEncoding.Encode(b, id[:])
	shortenedID := string(b)[:length]

	return shortenedID, nil
}

// GenerateKey 生成随机ID（小写+数字）
func GenerateKey(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("长度必须大于0")
	}

	letters := "abcdefghijkmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("生成随机数失败: %w", err)
	}

	for i, v := range b {
		b[i] = letters[v%byte(len(letters))]
	}
	return string(b), nil
}
