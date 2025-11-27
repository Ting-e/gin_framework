package router

import (
	"project/examples/gorm_crud/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有 API 路由
func ExampleRegisterRoutes_Grom(r *gin.Engine) {

	// 测试接口（gorm映射）
	gorm := r.Group("/gorm")
	{
		gorm.GET("/obtain-list", handler.GetList)     // 获取列表
		gorm.GET("/obtain-data/:id", handler.GetData) // 获取详情
		gorm.POST("/add-data", handler.AddData)       // 新增数据
		gorm.DELETE("/del-data/:id", handler.DelData) // 删除数据
		gorm.POST("/edit-data/:id", handler.EditData) // 修改数据
	}
}
