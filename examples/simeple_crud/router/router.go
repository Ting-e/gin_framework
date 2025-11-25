package router

import (
	"project/examples/simeple_crud/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有 API 路由
func ExampleRegisterRoutes(r *gin.Engine) {

	// 测试接口
	applet := r.Group("/api")
	{
		applet.GET("/obtain-list", handler.GetList)     // 获取列表
		applet.GET("/obtain-data/:id", handler.GetData) // 获取详情
		applet.POST("/add-data", handler.AddData)       // 新增数据
		applet.DELETE("/del-data/:id", handler.DelData) // 删除数据
		applet.POST("/edit-data/:id", handler.EditData) // 修改数据
	}
}
