package config

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"clean-archi/middleware"
	"clean-archi/routes"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New()

	// Middleware global
	app.Use(middleware.LoggerMiddleware)

	// Register semua routes
	routes.AuthRoutes(app, db)
	routes.AlumniRoutes(app, db)
	routes.JobRoutes(app, db)

	return app
}