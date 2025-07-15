package util

import "github.com/google/uuid"

// GenerateUUID 生成 UUID
func GenerateUUID() string {
	return uuid.New().String()
}
