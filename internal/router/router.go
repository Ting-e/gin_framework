package router

import (
	"project/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(r *gin.Engine) {
	// 测试接口
	applet := r.Group("/api")
	{
		applet.GET("/obtain-data", handler.GetData)
		applet.GET("/obtain-data/:id", handler.GetData)
		applet.POST("/obtain-data", handler.GetData)
		applet.DELETE("/obtain-data/:id", handler.GetData)
	}
}
