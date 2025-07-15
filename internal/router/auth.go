package router

import (
	"bmstock/internal/handler"
	"bmstock/pkg/di"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes 设置认证相关路由
func SetupAuthRoutes(api fiber.Router, container *di.Container) {
	var authHandler *handler.AuthHandler
	if err := container.Container.Invoke(func(h *handler.AuthHandler) { authHandler = h }); err != nil {
		panic(err)
	}

	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
}
