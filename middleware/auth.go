package middleware

import (
	"clean-archi/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware untuk semua endpoint yang butuh login
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token diperlukan",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Format token salah (gunakan 'Bearer <token>')",
			})
		}

		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		// Simpan data ke context untuk digunakan di handler berikutnya
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// Middleware khusus admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleValue := c.Locals("role")
		if roleValue == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User belum terautentikasi",
			})
		}

		role, ok := roleValue.(string)
		if !ok || role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Hanya admin yang diizinkan mengakses endpoint ini",
			})
		}

		return c.Next()
	}
}
