package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterDefaultMiddlewares 注册默认中间件
func RegisterDefaultMiddlewares(engine *gin.Engine, version string) {
	// 核心中间件（按顺序）
	engine.Use(Recovery())                // Panic 恢复
	engine.Use(Logger())                  // 日志记录
	engine.Use(CORS(version))             // 跨域处理
	engine.Use(Timeout(30 * time.Second)) // 请求超时
	// engine.Use(RateLimit())                  // 限流（可选）
}
