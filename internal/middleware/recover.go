package middleware

import (
	"bmstock/pkg/di"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Recover 全局 panic 捕获中间件
func Recover(container *di.Container) fiber.Handler {
	var logger *zap.Logger
	if err := container.Container.Invoke(func(l *zap.Logger) { logger = l }); err != nil {
		panic(err)
	}

	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// 记录 panic 日志
				logger.Error("Panic Recovered",
					zap.String("method", c.Method()),
					zap.String("path", c.Path()),
					zap.Any("query", c.Queries()),
					zap.String("body", string(c.Body())),
					zap.Any("headers", c.GetReqHeaders()),
					zap.Any("panic", r),
				)

				// 返回 500 错误响应
				err := c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":    fiber.StatusInternalServerError,
					"message": fmt.Sprintf("Internal Server Error: %v", r),
					"data":    nil,
				})
				if err != nil {
					logger.Error("Failed to send error response", zap.Error(err))
				}
			}
		}()

		return c.Next()
	}
}
