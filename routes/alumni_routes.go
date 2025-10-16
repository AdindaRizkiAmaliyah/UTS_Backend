package routes

import (
    "database/sql"
    "clean-archi/app/service"
    "github.com/gofiber/fiber/v2"
)

func AlumniRoutes(router fiber.Router, db *sql.DB) {
    // Inisialisasi service langsung tanpa handler
    alumniService := service.NewAlumniService(db)

    alumni := router.Group("/unair/alumni")
    alumni.Get("/", alumniService.GetAll, )
    alumni.Get("/:id", alumniService.GetByID)
    alumni.Post("/", alumniService.Create)
    alumni.Put("/:id", alumniService.Update)
    alumni.Delete("/:id", alumniService.Delete)
    
    // Tambahkan route untuk check alumni jika diperlukan
    alumni.Post("/check/:key", alumniService.CheckAlumni)
}