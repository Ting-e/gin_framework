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
		applet.GET("/obtain-list", handler.GetList)
		applet.GET("/obtain-data/:id", handler.GetData)
		applet.POST("/add-data", handler.AddData)
		applet.DELETE("/del-data/:id", handler.DelData)
		applet.POST("/edit-data/:id", handler.EditData)
	}
}
