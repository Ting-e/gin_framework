package middleware

import (
	"runtime/debug"

	"project/pkg/errcode"
	"project/pkg/logger"
	"project/pkg/response"

	"github.com/gin-gonic/gin"
)

// Recovery Panic 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := string(debug.Stack())
				logger.Sugar.Errorf("[Recovery] panic recovered:\n%v\n%s", err, stack)

				// 返回统一错误响应
				response.Failed(c, errcode.ServerError, errcode.ErrorMessage[errcode.ServerError])
				c.Abort()
			}
		}()
		c.Next()
	}
}
