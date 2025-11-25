package main

import (
	"project/pkg/logger"

	"project/internal/app"
	"project/internal/router"
)

func main() {

	application := app.MustInitApp()

	r := application.GetRouter()

	// 注册路由
	router.RegisterRoutes(r)

	if err := application.Run(); err != nil {
		logger.Sugar.Error(err)
	}
}
