package middleware

import (
	"strings"

	"project/pkg/jwt"
	"project/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 验证 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			// 根据错误类型返回不同消息
			var message string
			switch err {
			case jwt.ErrTokenExpired:
				message = "令牌已过期"
			case jwt.ErrTokenMalformed:
				message = "令牌格式错误"
			case jwt.ErrTokenNotValidYet:
				message = "令牌尚未生效"
			default:
				message = "令牌无效"
			}
			response.Unauthorized(c, message)
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// JWTAuthOptional 可选的 JWT 认证（不强制要求 token）
func JWTAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			if claims, err := jwt.ParseToken(tokenString); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
			}
		}

		c.Next()
	}
}

// RequireRole 角色权限验证中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "未找到用户角色信息")
			c.Abort()
			return
		}

		role := userRole.(string)
		hasPermission := false
		for _, r := range roles {
			if role == r {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (int64, bool) {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(int64); ok {
			return id, true
		}
	}
	return 0, false
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name, true
		}
	}
	return "", false
}

// GetRole 从上下文获取角色
func GetRole(c *gin.Context) (string, bool) {
	if role, exists := c.Get("role"); exists {
		if r, ok := role.(string); ok {
			return r, true
		}
	}
	return "", false
}
