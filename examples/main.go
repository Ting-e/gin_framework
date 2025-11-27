package main

import (
	"project/pkg/logger"

	"project/internal/app"

	gorm_crud_router "project/examples/gorm_crud/router"
	simeple_crud_router "project/examples/simeple_crud/router"
)

func main() {

	application := app.MustInitApp()

	r := application.GetRouter()

	// 注册路由
	simeple_crud_router.ExampleRegisterRoutes_SQL(r)
	gorm_crud_router.ExampleRegisterRoutes_Grom(r)

	if err := application.Run(); err != nil {
		logger.Sugar.Error(err)
	}
}
