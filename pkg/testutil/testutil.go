package testutil

import (
	"bmstock/api/v1"
	"bmstock/internal/model"
	"bmstock/pkg/util"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// NewFiberContext 创建 Fiber 测试上下文
func NewFiberContext(method, path string, body interface{}, headers map[string]string) (*fiber.Ctx, *fasthttp.RequestCtx) {
	// 创建 Fiber 应用
	app := fiber.New()

	// 创建 fasthttp 请求上下文
	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(method)
	fctx.Request.SetRequestURI(path)

	// 设置请求 body
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		fctx.Request.SetBody(bodyBytes)
		fctx.Request.Header.Set("Content-Type", "application/json")
	}

	// 设置请求头
	for k, v := range headers {
		fctx.Request.Header.Set(k, v)
	}

	// 创建 Fiber 上下文
	ctx := app.AcquireCtx(fctx)
	return ctx, fctx
}

// GenerateTestUser 生成测试用户
func GenerateTestUser(username, password string) model.User {
	hashedPassword, _ := util.GeneratePasswordHash(password)
	return model.User{
		ID:       util.GenerateUUID(),
		Username: username,
		Password: hashedPassword,
		Avatar:   "https://r2.example.com/avatar.jpg",
	}
}

// GenerateLoginRequest 生成登录请求
func GenerateLoginRequest(username, password string) v1.LoginRequest {
	return v1.LoginRequest{
		Username: username,
		Password: password,
	}
}

// GenerateChangePasswordRequest 生成修改密码请求
func GenerateChangePasswordRequest(oldPassword, newPassword string) v1.ChangePasswordRequest {
	return v1.ChangePasswordRequest{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
}

// GenerateRealNameAuthRequest 生成实名认证请求
func GenerateRealNameAuthRequest() v1.RealNameAuthRequest {
	return v1.RealNameAuthRequest{
		IDCardNumber: "123456789012345678",
		BankCard:     "1234567890123456",
		IDCardFront:  "https://r2.example.com/id_front.jpg",
		IDCardBack:   "https://r2.example.com/id_back.jpg",
	}
}
