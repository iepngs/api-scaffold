package config

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	App        AppConfig
	MySQL      MySQLConfig
	Redis      RedisConfig
	Cloudflare CloudflareConfig
	JWT        JWTConfig
}

// AppConfig 应用配置
type AppConfig struct {
	Port string
}

// MySQLConfig MySQL 配置
type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// DSN 生成 MySQL DSN
func (c MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// CloudflareConfig Cloudflare R2 配置
type CloudflareConfig struct {
	AccessKey string
	SecretKey string
	AccountID string
	Bucket    string
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string
}

// NewConfig 加载配置
func NewConfig() (*Config, error) {
	// 解析命令行参数
	configPath := flag.String("config", "./config.yaml", "path to config file")
	flag.Parse()

	// 初始化 Viper
	v := viper.New()
	v.SetConfigType("yaml")

	// 设置配置文件路径
	configFile, err := filepath.Abs(*configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}
	v.SetConfigFile(configFile)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}

	// 反序列化配置
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
