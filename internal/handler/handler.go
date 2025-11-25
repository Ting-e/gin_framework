package handler

import (
	srv "project/internal/service"
	"project/pkg/config"
	"time"

	"github.com/gin-gonic/gin"
)

var service srv.APIService

func init() {
	service = srv.GetService()
}

type PageData struct {
	Title         string // 页面标题
	ServiceName   string // 服务名称
	Now           string // 服务器时间
	Environment   string // 环境 (Development/Production)
	Version       string // 版本号
	StatusMessage string // 状态消息
}

func HomePage(c *gin.Context) {
	c.HTML(200, "index.html", PageData{
		Title:         "Gin_framework",
		ServiceName:   config.Get().Server.Name,
		Now:           time.Now().Format("2006-01-02 15:04:05"),
		Environment:   config.Get().Server.Environment,
		Version:       config.Get().Server.Version,
		StatusMessage: config.Get().Server.Name + " is running",
	})
}
