package router

import (
	"project/examples/simple_crud/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有 API 路由
func ExampleRegisterRoutes_SQL(r *gin.Engine) {

	// 测试接口（原生sql语句）
	sql := r.Group("/sql")
	{
		sql.GET("/obtain-list", handler.GetList)     // 获取列表
		sql.GET("/obtain-data/:id", handler.GetData) // 获取详情
		sql.POST("/add-data", handler.AddData)       // 新增数据
		sql.DELETE("/del-data/:id", handler.DelData) // 删除数据
		sql.POST("/edit-data/:id", handler.EditData) // 修改数据
	}

}
