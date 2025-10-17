package middleware

import (
	"clean-archi/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware untuk semua endpoint yang butuh login
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Jika tidak ada token, langsung kembalikan 401 Unauthorized
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token diperlukan",
			})
		}

		// Token biasanya dikirim dalam format "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			// Jika format salah, kembalikan error
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Format token salah (gunakan 'Bearer <token>')",
			})
		}

		// Validasi token menggunakan utilitas ValidateToken
		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			// Jika token tidak valid atau expired
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		// Simpan informasi user ke context (Locals)
		// Data ini akan dipakai di handler berikutnya (misal userID, email, role)
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		// Lanjut ke handler berikutnya
		return c.Next()
	}
}

// Middleware khusus admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil role dari Locals yang sebelumnya diset oleh AuthRequired
		roleValue := c.Locals("role")
		if roleValue == nil {
			// Jika tidak ada role, user belum login atau token invalid
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User belum terautentikasi",
			})
		}

		// Type assertion untuk memastikan role adalah string
		role, ok := roleValue.(string)
		if !ok || role != "admin" {
			// Jika bukan admin, kembalikan 403 Forbidden
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Hanya admin yang diizinkan mengakses endpoint ini",
			})
		}

		// Jika role admin, lanjut ke handler berikutnya
		return c.Next()
	}
}
