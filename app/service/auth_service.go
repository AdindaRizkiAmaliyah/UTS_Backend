package service

import (
	"clean-archi/app/model"
	"clean-archi/utils"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// Login menggunakan email dan password
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Body tidak valid",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email dan password wajib diisi",
		})
	}

	var alumni model.Alumni
	var hashedPassword string

	query := `
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, password, no_telepon, alamat, role, created_at, updated_at
		FROM alumni WHERE email = $1
	`
	err := s.db.QueryRow(query, req.Email).Scan(
		&alumni.ID,
		&alumni.NIM,
		&alumni.Nama,
		&alumni.Jurusan,
		&alumni.Angkatan,
		&alumni.TahunLulus,
		&alumni.Email,
		&hashedPassword,
		&alumni.NoTelp,
		&alumni.Alamat,
		&alumni.Role,
		&alumni.CreatedAt,
		&alumni.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Email atau password salah",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal query ke database",
			"detail": err.Error(), // tambahkan sementara untuk debug
		})
	}

	// Cek password
	if !utils.CheckPassword(req.Password, hashedPassword) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Password salah",
		})
	}

	// Generate JWT
	token, err := utils.GenerateToken(alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.LoginResponse{
		Alumni: alumni,
		Token:  token,
	})
}

// Profile dari token
func (s *AuthService) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	email := c.Locals("email").(string)
	role := c.Locals("role").(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}

// Handler opsional
func LoginHandler(c *fiber.Ctx, db *sql.DB) error {
	service := NewAuthService(db)
	return service.Login(c)
}

func ProfileHandler(c *fiber.Ctx, db *sql.DB) error {
	service := NewAuthService(db)
	return service.Profile(c)
}
