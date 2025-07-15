package router

import (
	"bmstock/internal/middleware"
	"bmstock/pkg/di"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes 初始化路由
func SetupRoutes(app *fiber.App, container *di.Container) {
	// 应用中间件
	app.Use(middleware.Recover(container))
	app.Use(middleware.CORS())
	app.Use(middleware.Logger(container))

	// API 版本前缀
	api := app.Group("/api/v1")

	// 公共路由
	SetupPublicRoutes(api, container)

	// 认证路由
	SetupAuthRoutes(api, container)

	// 用户路由
	SetupUserRoutes(api, container)
}
