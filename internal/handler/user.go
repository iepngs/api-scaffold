package handler

import (
	"bmstock/api/v1"
	"bmstock/internal/response"
	"bmstock/internal/service"

	"github.com/gofiber/fiber/v2"
)

// UserHandler 用户处理
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetProfile 获取个人资料
// @Summary 获取用户个人资料
// @Description 获取当前登录用户的个人资料
// @Tags 用户
// @Accept json
// @Produce json
// @Success 200 {object} v1.ProfileResponse
// @Failure 400 {object} response.Response
// @Router /api/v1/user/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string) // 从 JWT 中获取用户 ID
	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		return response.Error(c, 400, err.Error())
	}
	return response.Success(c, profile)
}

// ChangePassword 修改密码
// @Summary 修改用户密码
// @Description 用户修改密码，需提供旧密码和新密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body v1.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/user/change-password [post]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	var req v1.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, "Invalid request body")
	}

	userID := c.Locals("user_id").(string)
	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return response.Error(c, 400, err.Error())
	}
	return response.Success(c, nil)
}

// RealNameAuth 提交实名认证
// @Summary 提交实名认证
// @Description 用户提交身份证、银行卡等信息进行实名认证
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body v1.RealNameAuthRequest true "实名认证请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/user/realname-auth [post]
func (h *UserHandler) RealNameAuth(c *fiber.Ctx) error {
	var req v1.RealNameAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, "Invalid request body")
	}

	userID := c.Locals("user_id").(string)
	if err := h.userService.SubmitRealNameAuth(userID, req); err != nil {
		return response.Error(c, 400, err.Error())
	}
	return response.Success(c, nil)
}
