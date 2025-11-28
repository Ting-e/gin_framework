package main

import (
	"project/pkg/logger"

	"project/internal/app"

	gorm_crud_router "project/examples/gorm_crud/router"
	jwt_auth_router "project/examples/jwt_auth/router"
	simple_crud_router "project/examples/simple_crud/router"
)

func main() {

	application := app.MustInitApp()

	r := application.GetRouter()

	// 注册路由
	simple_crud_router.ExampleRegisterRoutes_SQL(r)
	gorm_crud_router.ExampleRegisterRoutes_Grom(r)
	jwt_auth_router.RegisterAuthRoutes(r)

	if err := application.Run(); err != nil {
		logger.Sugar.Error(err)
	}
}
