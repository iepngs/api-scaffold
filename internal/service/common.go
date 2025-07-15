package service

import (
	"bmstock/config"
	"bmstock/pkg/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CommonService 认证服务
type CommonService struct {
	db          *gorm.DB
	redisClient *redis.Client
	config      *config.Config
	logger      *zap.Logger
}

// NewCommonService 创建认证服务
func NewCommonService(db *gorm.DB, redisClient *redis.Client, config *config.Config, logger *zap.Logger) *CommonService {
	return &CommonService{db: db, redisClient: redisClient, config: config, logger: logger}
}
