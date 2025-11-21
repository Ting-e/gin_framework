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
	"project/pkg/logger"

	"project/internal/app"
	"project/internal/router"
)

func main() {

	application := app.MustInitApp()

	r := application.GetRouter()

	// 注册所有路由 ← 核心变化在这里
	router.RegisterRoutes(r)

	if err := application.Run(); err != nil {
		logger.Sugar.Error(err)
	}
}
