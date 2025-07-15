package service

import (
	"bmstock/config"
	"bmstock/pkg/redis"
	"bmstock/pkg/testutil"
	"bmstock/pkg/util"
	"errors"
	"github.com/go-redis/redismock/v9"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAuthService_Login(t *testing.T) {
	// Create Redis mock
	dbRedis, mockRedis := redismock.NewClientMock()
	redisClient := &redis.Client{Client: dbRedis}

	// Create SQL mock
	dbSQL, mockSQL, err := sqlmock.New()
	assert.NoError(t, err)

	// Initialize GORM with sqlmock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: dbSQL}), &gorm.Config{})
	assert.NoError(t, err)

	// Mock GORM's initialization query (SELECT VERSION())
	mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT VERSION()")).WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("8.0.27"))

	// Config
	cfg := &config.Config{JWT: config.JWTConfig{Secret: "test-secret"}}

	// Create service
	s := NewAuthService(gormDB, redisClient, cfg)

	// Simulate user
	hashedPassword, _ := util.GeneratePasswordHash("pass1234")
	user := testutil.GenerateTestUser("user123", "pass1234")

	// Test cases
	tests := []struct {
		name        string
		username    string
		password    string
		setupMock   func()
		expectToken bool
		expectErr   bool
	}{
		{
			name:     "Valid credentials",
			username: "user123",
			password: "pass1234",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "avatar"}).
					AddRow(user.ID, user.Username, hashedPassword, user.Avatar)
				mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
					WithArgs("user123").WillReturnRows(rows)
				mockRedis.ExpectSet("token:"+user.ID, "", 24*time.Hour).SetVal("OK")
			},
			expectToken: true,
			expectErr:   false,
		},
		{
			name:     "User not found",
			username: "user123",
			password: "pass1234",
			setupMock: func() {
				mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
					WithArgs("user123").WillReturnError(gorm.ErrRecordNotFound)
			},
			expectToken: false,
			expectErr:   true,
		},
		{
			name:     "Invalid password",
			username: "user123",
			password: "wrongpass",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "avatar"}).
					AddRow(user.ID, user.Username, hashedPassword, user.Avatar)
				mockSQL.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
					WithArgs("user123").WillReturnRows(rows)
			},
			expectToken: false,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			token, _, err := s.Login(tt.username, tt.password)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}

	// Verify mock expectations
	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.NoError(t, mockRedis.ExpectationsWereMet())
}

func TestAuthService_Logout(t *testing.T) {
	// Create Redis mock
	db, mock := redismock.NewClientMock()
	redisClient := &redis.Client{Client: db}

	// Create SQL mock (not used in Logout)
	dbSQL, _, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: dbSQL}), &gorm.Config{})
	assert.NoError(t, err)

	// Config
	cfg := &config.Config{JWT: config.JWTConfig{Secret: "test-secret"}}

	// Create service
	s := NewAuthService(gormDB, redisClient, cfg)

	// Test cases
	t.Run("Successful logout", func(t *testing.T) {
		mock.ExpectDel("token:test-token").SetVal(1)
		err := s.Logout("test-token")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed logout", func(t *testing.T) {
		mock.ExpectDel("token:test-token").SetErr(errors.New("redis error"))
		err := s.Logout("test-token")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
