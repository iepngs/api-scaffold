package handler

import (
	"bmstock/api/v1"
	"bmstock/internal/response"
	"bmstock/internal/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler 认证处理
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 处理登录请求
// @Summary 用户登录
// @Description 用户使用用户名和密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "登录请求"
// @Success 200 {object} v1.LoginResponse
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req v1.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, "Invalid request body")
	}

	token, expire, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, err.Error())
	}

	return response.Success(c, v1.LoginResponse{Token: token, Expire: expire})
}

// Logout 处理退出请求
// @Summary 用户退出
// @Description 用户退出登录
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} v1.LogoutResponse
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// 假设使用 JWT，退出可清理 Redis 中的 token
	if err := h.authService.Logout(c.Get("Authorization")); err != nil {
		return response.Error(c, 400, err.Error())
	}

	return response.Success(c, v1.LogoutResponse{Message: "Logged out successfully"})
}
