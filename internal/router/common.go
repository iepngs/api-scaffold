package router

import (
	"bmstock/internal/handler"
	"bmstock/pkg/di"

	"github.com/gofiber/fiber/v2"
)

// SetupPublicRoutes 设置公共路由
func SetupPublicRoutes(api fiber.Router, container *di.Container) {
	var commonHandler *handler.CommonHandler
	if err := container.Container.Invoke(func(h *handler.CommonHandler) { commonHandler = h }); err != nil {
		panic(err)
	}

	auth := api.Group("/public")
	auth.Post("/deploy", commonHandler.Deploy)
}
