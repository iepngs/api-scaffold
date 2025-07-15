package middleware

import (
	"bmstock/pkg/di"
	"bmstock/pkg/redis"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// Auth JWT 认证中间件
func Auth(container *di.Container) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
				"code":    fiber.StatusUnauthorized,
				"message": "Missing token",
			})
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(container.Config.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
				"code":    fiber.StatusUnauthorized,
				"message": "Invalid token",
			})
		}

		// 验证 Redis 中 token 是否存在
		var redisClient *redis.Client
		if err = container.Invoke(func(r *redis.Client) { redisClient = r }); err != nil {
			return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
				"code":    fiber.StatusInternalServerError,
				"message": "Internal server error",
			})
		}
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["sub"].(int64)
		if _, err = redisClient.Get(c.Context(), fmt.Sprintf("token:%d", userID)).Result(); err != nil {
			return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
				"code":    fiber.StatusUnauthorized,
				"message": "Token expired or invalid",
			})
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
