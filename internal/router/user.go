package router

import (
	"bmstock/internal/handler"
	"bmstock/internal/middleware"
	"bmstock/pkg/di"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes 设置用户相关路由
func SetupUserRoutes(api fiber.Router, container *di.Container) {
	var userHandler *handler.UserHandler
	if err := container.Container.Invoke(func(h *handler.UserHandler) { userHandler = h }); err != nil {
		panic(err)
	}

	user := api.Group("/user")
	user.Use(middleware.Auth(container)) // 使用 *di.Container
	user.Get("/profile", userHandler.GetProfile)
	user.Post("/change-password", userHandler.ChangePassword)
	user.Post("/realname-auth", userHandler.RealNameAuth)
}
