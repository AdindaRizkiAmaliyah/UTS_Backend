package routes

import (
	"clean-archi/app/service"
	"clean-archi/middleware"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func JobRoutes(router fiber.Router, db *sql.DB) {
	jobService := service.NewJobService(db)

	// Kelompok route pekerjaan
	job := router.Group("/unair/pekerjaan")

	// === [PUBLIC ROUTES] ===
	job.Get("/", jobService.GetAll)
	job.Get("/:id", jobService.GetByID)

	// === [PROTECTED ROUTES - butuh JWT valid] ===
	protected := job.Group("", middleware.AuthRequired())

	// Pekerjaan milik alumni tertentu
	protected.Get("/alumni/:alumni_id", jobService.GetJobsByAlumniID)

	// Tambah pekerjaan (harus login)
	protected.Post("/", jobService.Create)

	// Update pekerjaan sendiri
	protected.Put("/:id", jobService.Update)

	// Soft delete pekerjaan sendiri (harus login)
	protected.Delete("/soft/:id", jobService.SoftDeleteByAlumni)

	// === [ADMIN ROUTES] ===
	admin := protected.Group("/admin", middleware.AdminOnly())

	// Soft delete semua pekerjaan milik alumni
	admin.Delete("/soft/:alumni_id", jobService.SoftDeleteAllByAdmin)

	// Lihat data yang dihapus (trash)
	admin.Get("/trashed", jobService.GetTrashed)

	// Restore pekerjaan yang dihapus
	admin.Put("/restore/:id", jobService.Restore)

	// Hard delete pekerjaan (hapus permanen)
	admin.Delete("/hard/:id", jobService.HardDelete)
}
