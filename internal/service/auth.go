package service

import (
	"bmstock/config"
	"bmstock/internal/model"
	"bmstock/pkg/redis"
	"bmstock/pkg/util"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthService 认证服务
type AuthService struct {
	db          *gorm.DB
	redisClient *redis.Client
	config      *config.Config
	logger      *zap.Logger
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, redisClient *redis.Client, config *config.Config, logger *zap.Logger) *AuthService {
	return &AuthService{db: db, redisClient: redisClient, config: config, logger: logger}
}

// Login 登录逻辑
func (s *AuthService) Login(username, password string) (string, int64, error) {
	// 查询用户
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, errors.New("user not found")
		}
		return "", 0, err
	}

	// 验证密码
	if !util.CheckPasswordHash(password, user.Password) {
		return "", 0, errors.New("invalid credentials")
	}

	duration := 24 * time.Hour
	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(duration).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", 0, err
	}

	// 存储 token 到 Redis
	err = s.redisClient.Set(context.Background(), fmt.Sprintf("token:%d", user.ID), tokenString, duration).Err()
	if err != nil {
		s.logger.Error("fail to set user token to redis", zap.Error(err))
	}

	return tokenString, time.Now().Add(duration).Unix(), nil
}

// Logout 退出逻辑
func (s *AuthService) Logout(token string) error {
	// 清理 Redis 中的 token
	return s.redisClient.Del(context.Background(), "token:"+token).Err()
}
