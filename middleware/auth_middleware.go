package middleware

import (
	"clean-archi/app/model"
	"clean-archi/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// 🔐 Middleware untuk semua endpoint yang butuh login
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
				"error": "Format token salah. Gunakan format 'Bearer <token>'",
			})
		}

		token := tokenParts[1]

		// Validasi token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		// Simpan klaim ke context
		c.Locals("user", claims)
		c.Locals("role", claims.Role)

		// ✅ Tambahkan ini supaya SoftDelete dan Restore bisa tahu siapa user-nya
		c.Locals("user_id", claims.UserID)

		return c.Next()
	}
}

// 👑 Middleware khusus admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userData := c.Locals("user")
		if userData == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User belum login atau token tidak valid",
			})
		}

		var role string
		switch v := userData.(type) {
		case *model.JWTClaims:
			role = v.Role
		case map[string]interface{}:
			if r, ok := v["role"].(string); ok {
				role = r
			}
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Format klaim token tidak dikenali",
			})
		}

		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Hanya admin yang diizinkan mengakses endpoint ini",
			})
		}

		return c.Next()
	}
}
