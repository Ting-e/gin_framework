/*
 * @Author: 廷益 15137125570@163.com
 * @Date: 2025-05-09 15:18:13
 * @LastEditors: 廷益 15137125570@163.com
 * @LastEditTime: 2025-11-21 11:55:00
 * @FilePath: /cmd/api/main.go
 * @Description: API 服务主入口
 */
package main

import (
	"flag"
	"project/internal/handler"
	"project/pkg/logger"

	"project/internal/app"

	"github.com/robfig/cron/v3"
)

var (
	application app.App
	conf        string
	logPath     string
)

func setupLogRotation() {
	c := cron.New()
	// 每天凌晨1点分割日志（注意 cron 表达式格式）
	c.AddFunc("0 0 1 * * *", func() {
		logger.InitLogger(logPath)
	})
	c.Start()
}

func init() {
	flag.StringVar(&conf, "conf", "./configs/app-dev.yaml", "配置文件路径")
	flag.StringVar(&logPath, "log", "./log/", "日志目录路径")
	flag.Parse()

	application = app.MustInitApp(conf, logPath)
	setupLogRotation()
}

func main() {
	r := application.GetRouter() // 假设路由由 pkg 管理

	login := r.Group("/applet/ceshi")
	login.GET("/obtain-data", handler.GetData)

	if err := application.Run(); err != nil {
		logger.Sugar.Error(err)
	}
}
