package main

import (
	"clean-archi/config"
	_ "clean-archi/docs" // WAJIB: import hasil `swag init`
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// ====================================================================
// ========================= SWAGGER CONFIG ===========================
// ====================================================================

// @title           Clean Archi API
// @version         1.0
// @description     API untuk pengelolaan data alumni menggunakan Clean Architecture dan MongoDB
// @termsOfService  https://github.com/yourname/clean-archi

// @contact.name   API Support
// @contact.url    https://github.com/yourname/clean-archi/issues
// @contact.email  support@cleanarchi.com

// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT

// @host      localhost:3000
// @BasePath  /api/v1
// @schemes   http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT dengan format: `Bearer <token>`


// ====================================================================
// =========================== MAIN FUNCTION ===========================
// ====================================================================

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Gagal memuat .env, lanjutkan dengan default...")
	}

	// Buat Fiber app manual agar bisa tambahkan middleware global
	app := fiber.New()

	// 🛠️ Tambahkan middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Swagger UI dan API sama domain
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Tambahkan semua route dari config.NewApp()
	internalApp := config.NewApp()
	app.Mount("/", internalApp)

	// Tambahkan Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Port default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Jalankan server di localhost:3000
	fmt.Printf("\n🚀 Server berjalan di http://localhost:%s\n", port)
	fmt.Printf("📘 Swagger Docs: http://localhost:%s/swagger/index.html\n", port)
	fmt.Println("🔑 Gunakan endpoint login JWT: POST /api/v1/auth/login")

	if err := app.Listen("localhost:" + port); err != nil {
		log.Fatal("❌ Gagal menjalankan server:", err)
	}
}
