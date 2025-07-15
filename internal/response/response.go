package response

import "github.com/gofiber/fiber/v2"

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code" example:"200"`   // 状态码
	Message string      `json:"message" example:"OK"` // 消息
	Data    interface{} `json:"data"`                 // 数据
}

// Success 返回成功响应
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Code:    200,
		Message: "OK",
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
