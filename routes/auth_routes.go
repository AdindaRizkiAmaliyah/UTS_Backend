package routes

import (
	"clean-archi/app/service"
	"clean-archi/middleware"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, db *sql.DB) {
	// Inisialisasi auth service
	authService := service.NewAuthService(db)
	
	api := app.Group("/api")

	// Public routes
	api.Post("/login", authService.Login)

	// Protected routes
	protected := api.Group("", middleware.AuthRequired())
	protected.Get("/profile", authService.Profile)
}