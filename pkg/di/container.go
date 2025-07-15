package di

import (
	"bmstock/config"
	"bmstock/internal/handler"
	"bmstock/internal/model"
	"bmstock/internal/service"
	"bmstock/pkg/cloudflare"
	"bmstock/pkg/logger"
	"bmstock/pkg/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.uber.org/dig"
)

// Container 依赖注入容器
type Container struct {
	*dig.Container
	Config *config.Config
}

// NewContainer 创建依赖注入容器
func NewContainer() (*Container, error) {
	c := dig.New()

	// 注册配置
	if err := c.Provide(config.NewConfig); err != nil {
		return nil, err
	}

	// 注册数据库
	if err := c.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		return model.InitDB(cfg)
	}); err != nil {
		return nil, err
	}

	// 注册 Redis
	if err := c.Provide(func(cfg *config.Config) (*redis.Client, error) {
		return redis.NewClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	}); err != nil {
		return nil, err
	}

	// 注册 Cloudflare R2
	if err := c.Provide(func(cfg *config.Config) *cloudflare.R2Client {
		return cloudflare.NewR2Client(
			cfg.Cloudflare.AccessKey,
			cfg.Cloudflare.SecretKey,
			cfg.Cloudflare.AccountID,
			cfg.Cloudflare.Bucket,
		)
	}); err != nil {
		return nil, err
	}

	// 注册 Zap 日志
	if err := c.Provide(logger.NewLogger); err != nil {
		return nil, err
	}

	// 注册服务
	if err := c.Provide(func(db *gorm.DB, redis *redis.Client, cfg *config.Config, logger *zap.Logger) *service.CommonService {
		return service.NewCommonService(db, redis, cfg, logger)
	}); err != nil {
		return nil, err
	}
	if err := c.Provide(func(db *gorm.DB, redis *redis.Client, cfg *config.Config, logger *zap.Logger) *service.AuthService {
		return service.NewAuthService(db, redis, cfg, logger)
	}); err != nil {
		return nil, err
	}
	if err := c.Provide(func(db *gorm.DB, r2 *cloudflare.R2Client) *service.UserService {
		return service.NewUserService(db, r2)
	}); err != nil {
		return nil, err
	}

	// 注册处理器
	if err := c.Provide(handler.NewCommonHandler); err != nil {
		return nil, err
	}
	if err := c.Provide(handler.NewAuthHandler); err != nil {
		return nil, err
	}
	if err := c.Provide(handler.NewUserHandler); err != nil {
		return nil, err
	}

	// 获取配置
	var cfg *config.Config
	if err := c.Invoke(func(c *config.Config) { cfg = c }); err != nil {
		return nil, err
	}

	return &Container{Container: c, Config: cfg}, nil
}
