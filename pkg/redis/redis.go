package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

// Client Redis 客户端
type Client struct {
	*redis.Client
}

// NewClient 创建 Redis 客户端
func NewClient(addr, password string, db int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}
