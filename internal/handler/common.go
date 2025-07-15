package handler

import (
	"bmstock/internal/response"
	"bmstock/internal/service"
	"bmstock/pkg/constant"
	"bmstock/pkg/util"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CommonHandler 认证处理
type CommonHandler struct {
	authService *service.AuthService
}

// NewCommonHandler 创建认证处理器
func NewCommonHandler(authService *service.AuthService) *CommonHandler {
	return &CommonHandler{authService: authService}
}

// Deploy 接收本地上传的编译包到服务器的 /tmp 目录
// 本地执行 make deploy
// /api/v1/public/deploy [post]
func (h *CommonHandler) Deploy(c *fiber.Ctx) error {
	token := c.Get("token")
	if token != "" && util.Md5Hash(token) == util.Md5Hash(time.Now().Format(constant.UploadBinaryAuthToken)) {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "无法获取上传的文件")
	}

	uploadDir := "/tmp"
	if err = os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "无法创建上传目录")
	}

	dst := filepath.Join(uploadDir, file.Filename)
	if err = c.SaveFile(file, dst); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "无法保存文件")
	}

	// 使用 bash 来执行组合命令
	cmd := exec.Command("/bin/bash", "-c", `nohup sleep 5; sh /data/wwwroot/api/deploy.sh > /dev/null 2>&1 &`)
	if err = cmd.Start(); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.Map{
		"filename": file.Filename,
		"path":     dst,
	})
}
