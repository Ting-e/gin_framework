package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Code    int    // 状态码
	Message string // 消息
}

// PostRequest 发送post请求
//
// 发送post请求
//
//	url：请求路径
//	requestBody：请求数据
//	contentType：请求数据格式（默认：application/json）
func PostRequest(url string, requestBody interface{}, contentType string) ([]byte, error) {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("无法封送请求正文: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("无法创建HTTP请求: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法执行HTTP请求: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 || err != nil {

		return nil, fmt.Errorf("收到未成功的状态码: %d, 响应体: %s", resp.StatusCode, string(bodyBytes))
	}

	return bodyBytes, nil
}
