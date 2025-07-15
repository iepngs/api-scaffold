package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	ID           int64  `gorm:"type:bigint;primaryKey"`           // 用户 ID (UUID)
	Username     string `gorm:"type:varchar(50);unique;not null"` // 用户名
	Password     string `gorm:"type:varchar(100);not null"`       // 密码（bcrypt 加密）
	Avatar       string `gorm:"type:varchar(255)"`                // 头像 URL
	IDCardNumber string `gorm:"type:varchar(18)"`                 // 身份证号
	BankCard     string `gorm:"type:varchar(20)"`                 // 银行卡号
	IDCardFront  string `gorm:"type:varchar(255)"`                // 身份证正面 URL
	IDCardBack   string `gorm:"type:varchar(255)"`                // 身份证反面 URL
	IsVerified   bool   `gorm:"default:false"`                    // 是否实名认证
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"` // 软删除
}
