package service

import (
	"clean-archi/app/model"
	"clean-archi/app/repository/MongoRepo"
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ============================================================
// ================          STRUCT          ==================
// ============================================================

// AuthService bertanggung jawab untuk proses registrasi dan login alumni
type AuthService struct {
	mongoRepo *MongoRepo.AlumniMongoRepository
}

// NewAuthService membuat instance baru dari AuthService
func NewAuthService(mongoRepo *MongoRepo.AlumniMongoRepository) *AuthService {
	return &AuthService{mongoRepo: mongoRepo}
}

// ============================================================
// ================          REGISTER         =================
// ============================================================

// Register godoc
// @Summary Registrasi akun alumni baru
// @Description Membuat akun baru untuk alumni, dengan enkripsi password sebelum disimpan ke database
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.Alumni true "Data Alumni untuk registrasi"
// @Success 201 {object} model.Alumni
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (s *AuthService) Register(c *fiber.Ctx) error {
	var req model.Alumni
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Body request tidak valid",
			"success": false,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cek apakah email sudah terdaftar
	count, err := s.mongoRepo.GetCollection().CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memeriksa email: " + err.Error(),
			"success": false,
		})
	}

	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email sudah terdaftar",
			"success": false,
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengenkripsi password",
			"success": false,
		})
	}

	req.Password = string(hashedPassword)
	req.Role = "user"
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	alumni, err := s.mongoRepo.Create(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal registrasi: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registrasi berhasil",
		"success": true,
		"data":    alumni,
	})
}

// ============================================================
// ================            LOGIN          =================
// ============================================================

// Login godoc
// @Summary Login alumni menggunakan email dan password
// @Description Mengecek kredensial user (email dan password), lalu mengembalikan JWT token jika valid
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "Data login (email dan password)"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Body request tidak valid",
			"success": false,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var alumni model.Alumni
	err := s.mongoRepo.GetCollection().FindOne(ctx, bson.M{"email": req.Email}).Decode(&alumni)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Email tidak ditemukan",
			"success": false,
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mencari email: " + err.Error(),
			"success": false,
		})
	}

	// Bandingkan password
	err = bcrypt.CompareHashAndPassword([]byte(alumni.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Password salah",
			"success": false,
		})
	}

	// Buat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		UserID: alumni.MongoID.Hex(),
		Email:  alumni.Email,
		Role:   alumni.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "unair-auth-service",
		},
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat token JWT",
			"success": false,
		})
	}

	response := model.LoginResponse{
		Alumni: alumni,
		Token:  tokenString,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login berhasil",
		"success": true,
		"data":    response,
	})
}
