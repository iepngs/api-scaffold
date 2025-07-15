package model

import (
	"bmstock/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.MySQL.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	return db, nil
}
