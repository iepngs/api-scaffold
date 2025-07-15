package database

import (
	"bmstock/config"
	"bmstock/internal/model"
	"gorm.io/gorm"
)

// NewMySQL 初始化 MySQL 数据库
func NewMySQL(cfg *config.Config) (*gorm.DB, error) {
	return model.InitDB(cfg)
}
