package router

import (
	"project/examples/jwt_auth/handler"
	"project/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(r *gin.Engine) {
	// 公开路由（不需要认证）
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", handler.Register)    // 注册
		authGroup.POST("/login", handler.Login)          // 登录
		authGroup.POST("/refresh", handler.RefreshToken) // 刷新 token
		authGroup.POST("/logout", handler.Logout)        // 登出
	}

	// 需要认证的路由
	protectedGroup := r.Group("/api/auth")
	protectedGroup.Use(middleware.JWTAuth()) // 使用 JWT 中间件
	{
		protectedGroup.GET("/userinfo", handler.GetUserInfo)            // 获取用户信息
		protectedGroup.POST("/change-password", handler.ChangePassword) // 修改密码
		protectedGroup.POST("/logout-all", handler.LogoutAllDevices)    // 登出所有设备
	}
}
