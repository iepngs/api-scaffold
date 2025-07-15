package util

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Hash 计算给定字符串的MD5哈希值并返回其十六进制表示。
func Md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
