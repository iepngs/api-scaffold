package service

import (
	"bmstock/api/v1"
	"bmstock/pkg/cloudflare"
	"bmstock/pkg/util"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserService_GetProfile(t *testing.T) {
	// 创建 SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	// 创建 R2 mock（简化）
	r2Client := &cloudflare.R2Client{}

	// 创建服务
	s := NewUserService(gormDB, r2Client)

	// 测试用例
	tests := []struct {
		name          string
		userID        string
		setupMock     func()
		expectErr     bool
		expectProfile v1.ProfileResponse
	}{
		{
			name:   "User found",
			userID: "uuid123",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "avatar"}).
					AddRow("uuid123", "user123", "https://r2.example.com/avatar.jpg")
				mock.ExpectQuery("SELECT").WithArgs("uuid123").WillReturnRows(rows)
			},
			expectErr: false,
			expectProfile: v1.ProfileResponse{
				ID:       "uuid123",
				Username: "user123",
				Avatar:   "https://r2.example.com/avatar.jpg",
			},
		},
		{
			name:   "User not found",
			userID: "uuid404",
			setupMock: func() {
				mock.ExpectQuery("SELECT").WithArgs("uuid404").WillReturnError(gorm.ErrRecordNotFound)
			},
			expectErr:     true,
			expectProfile: v1.ProfileResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			profile, err := s.GetProfile(tt.userID)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectProfile, profile)
			}
		})
	}

	// 验证 SQL mock
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_ChangePassword(t *testing.T) {
	// 创建 SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	// 创建 R2 mock（简化）
	r2Client := &cloudflare.R2Client{}

	// 创建服务
	s := NewUserService(gormDB, r2Client)

	// 模拟用户
	hashedPassword, _ := util.GeneratePasswordHash("oldpass123")

	// 测试用例
	tests := []struct {
		name        string
		userID      string
		oldPassword string
		newPassword string
		setupMock   func()
		expectErr   bool
	}{
		{
			name:        "Valid password change",
			userID:      "uuid123",
			oldPassword: "oldpass123",
			newPassword: "newpass123",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "password"}).
					AddRow("uuid123", hashedPassword)
				mock.ExpectQuery("SELECT").WithArgs("uuid123").WillReturnRows(rows)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:        "Invalid old password",
			userID:      "uuid123",
			oldPassword: "wrongpass",
			newPassword: "newpass123",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "password"}).
					AddRow("uuid123", hashedPassword)
				mock.ExpectQuery("SELECT").WithArgs("uuid123").WillReturnRows(rows)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := s.ChangePassword(tt.userID, tt.oldPassword, tt.newPassword)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// 验证 SQL mock
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserService_SubmitRealNameAuth(t *testing.T) {
	// 创建 SQL mock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	assert.NoError(t, err)

	// 创建 R2 mock（简化）
	r2Client := &cloudflare.R2Client{}

	// 创建服务
	s := NewUserService(gormDB, r2Client)

	// 测试用例
	req := v1.RealNameAuthRequest{
		IDCardNumber: "123456789012345678",
		BankCard:     "1234567890123456",
		IDCardFront:  "https://r2.example.com/id_front.jpg",
		IDCardBack:   "https://r2.example.com/id_back.jpg",
	}

	tests := []struct {
		name      string
		userID    string
		setupMock func()
		expectErr bool
	}{
		{
			name:   "Successful real name auth",
			userID: "uuid123",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow("uuid123")
				mock.ExpectQuery("SELECT").WithArgs("uuid123").WillReturnRows(rows)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:   "User not found",
			userID: "uuid404",
			setupMock: func() {
				mock.ExpectQuery("SELECT").WithArgs("uuid404").WillReturnError(gorm.ErrRecordNotFound)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := s.SubmitRealNameAuth(tt.userID, req)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// 验证 SQL mock
	assert.NoError(t, mock.ExpectationsWereMet())
}
