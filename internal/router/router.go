package router

import (
	"project/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(r *gin.Engine) {

	r.LoadHTMLFiles("web/index.html")

	// localhost
	loclahost := r.Group("/")
	{
		loclahost.GET("", handler.HomePage)
	}

}
