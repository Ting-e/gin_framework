package middleware

import (
	"sync"
	"time"

	"project/pkg/errcode"
	"project/pkg/response"

	"github.com/gin-gonic/gin"
)

// 简单的令牌桶限流器
type rateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func newRateLimiter(maxTokens int, refillRate time.Duration) *rateLimiter {
	return &rateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (r *rateLimiter) allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefill)

	// 补充令牌
	tokensToAdd := int(elapsed / r.refillRate)
	if tokensToAdd > 0 {
		r.tokens = min(r.tokens+tokensToAdd, r.maxTokens)
		r.lastRefill = now
	}

	// 检查是否有可用令牌
	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var globalLimiter = newRateLimiter(100, time.Second) // 每秒 100 个请求

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !globalLimiter.allow() {
			response.Failed(c, errcode.TooManyRequest, errcode.ErrorMessage[errcode.TooManyRequest])
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitWithConfig 自定义配置的限流中间件
func RateLimitWithConfig(maxTokens int, refillRate time.Duration) gin.HandlerFunc {
	limiter := newRateLimiter(maxTokens, refillRate)
	return func(c *gin.Context) {
		if !limiter.allow() {
			response.Failed(c, errcode.TooManyRequest, errcode.ErrorMessage[errcode.TooManyRequest])
			c.Abort()
			return
		}
		c.Next()
	}
}
