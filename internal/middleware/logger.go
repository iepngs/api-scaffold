package middleware

import (
	"bmstock/pkg/di"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger(container *di.Container) fiber.Handler {
	var logger *zap.Logger
	if err := container.Container.Invoke(func(l *zap.Logger) { logger = l }); err != nil {
		panic(err)
	}

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// 捕获请求 body
		body := c.Body()
		// 截断 body（最大 1KB）
		if len(body) > 1024 {
			body = body[:1024]
		}

		// 执行请求
		err := c.Next()
		duration := time.Since(start)

		// 获取响应数据
		responseData := c.Response().Body()
		// 截断 response（最大 1KB）
		if len(responseData) > 1024 {
			responseData = responseData[:1024]
		}

		// 基本日志字段
		logFields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		}

		// 错误请求（状态码 >= 400）记录详细信息
		if c.Response().StatusCode() >= 400 {
			logFields = append(logFields,
				zap.Any("query", c.Queries()),
				zap.String("body", string(body)),
				zap.Any("headers", c.GetReqHeaders()),
				zap.String("response", string(responseData)),
			)
			logger.Error("HTTP Request Error", logFields...)
		} else {
			// 正常请求仅记录基本信息
			logger.Info("HTTP Request", logFields...)
		}

		return err
	}
}
