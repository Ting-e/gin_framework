package middleware

import (
	"context"
	"net/http"
	"time"

	"project/pkg/response"

	"github.com/gin-gonic/gin"
)

// Timeout 超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		// 使用 channel 来检测处理是否完成
		finished := make(chan struct{})
		go func() {
			c.Next()
			close(finished)
		}()

		select {
		case <-finished:
			// 请求正常完成
			return
		case <-ctx.Done():
			// 请求超时
			if ctx.Err() == context.DeadlineExceeded {
				response.Failed(c, http.StatusRequestTimeout, "Request timeout")
				c.Abort()
			}
		}
	}
}
