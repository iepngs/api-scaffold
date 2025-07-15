package service

import (
	"bmstock/api/v1"
	"bmstock/internal/model"
	"bmstock/pkg/cloudflare"
	"bmstock/pkg/util"
	"errors"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db       *gorm.DB
	r2Client *cloudflare.R2Client
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB, r2Client *cloudflare.R2Client) *UserService {
	return &UserService{db: db, r2Client: r2Client}
}

// GetProfile 获取个人资料
func (s *UserService) GetProfile(userID string) (v1.ProfileResponse, error) {
	var user model.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return v1.ProfileResponse{}, errors.New("user not found")
	}
	return v1.ProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID, oldPassword, newPassword string) error {
	var user model.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if !util.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("invalid old password")
	}

	hashedPassword, err := util.GeneratePasswordHash(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.db.Save(&user).Error
}

// SubmitRealNameAuth 提交实名认证
func (s *UserService) SubmitRealNameAuth(userID string, req v1.RealNameAuthRequest) error {
	var user model.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	// 假设身份证图片已通过 R2 上传，存储 URL
	user.IDCardNumber = req.IDCardNumber
	user.BankCard = req.BankCard
	user.IDCardFront = req.IDCardFront
	user.IDCardBack = req.IDCardBack
	user.IsVerified = true
	return s.db.Save(&user).Error
}
