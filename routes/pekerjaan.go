package routes

import (
	"clean-archi/app/service" // Path service Anda
	"clean-archi/middleware" // Path middleware Anda
	"database/sql"
	// strconv sudah dipindahkan ke service/middleware

	"github.com/gofiber/fiber/v2"
)

// JobRoutes mengatur semua rute yang terkait dengan entitas Pekerjaan
func JobRoutes(router fiber.Router, db *sql.DB) {
	jobService := service.NewJobService(db)

	// Base Group: /unair/pekerjaan
	job := router.Group("/unair/pekerjaan")

	// === [2. PROTECTED ROUTES - Alumni/User (Membutuhkan AuthRequired)] ===
	protected := job.Group("", middleware.AuthRequired())

	// **[RUTE KHUSUS ALUMNI DITEMPATKAN DI AWAL GROUP PROTECTED]**
	// VIEW Pekerjaan yang di-soft delete milik sendiri (GET /unair/pekerjaan/trashed)
	protected.Get("/trashed", jobService.GetTrashed) 

	// GET Daftar pekerjaan milik sendiri (GET /unair/pekerjaan/my)
	// ID alumni diambil langsung dari token di service.GetMyJobs.
	protected.Get("/my", jobService.GetMyJobs) 

	// RESTORE Pekerjaan milik sendiri (PUT /unair/pekerjaan/restore/:id)
	protected.Put("/restore/:id", jobService.Restore)

	// === [3. ADMIN ROUTES - Membutuhkan AuthRequired + AdminOnly] ===
	admin := job.Group("/admin", middleware.AuthRequired(), middleware.AdminOnly())

	// Admin: Melihat SEMUA data pekerjaan yang dihapus (GET /unair/pekerjaan/admin/trashed/all)
	admin.Get("/trashed/all", jobService.GetTrashed) 

	// Admin: Melihat pekerjaan alumni tertentu (GET /unair/pekerjaan/admin/alumni/:alumni_id)
	admin.Get("/alumni/:alumni_id", jobService.GetJobsByAlumniID)

	// Admin: Restore pekerjaan siapa pun (PUT /unair/pekerjaan/admin/restore/:id)
	admin.Put("/restore/:id", jobService.Restore)

	// Admin: Soft delete SEMUA pekerjaan milik alumni tertentu (DELETE /unair/pekerjaan/admin/soft/alumni/:alumni_id)
	admin.Delete("/soft/alumni/:alumni_id", jobService.SoftDeleteAllByAdmin)

	// Admin: Hard delete pekerjaan siapa pun (DELETE /unair/pekerjaan/admin/hard/:id)
	admin.Delete("/hard/:id", jobService.HardDelete)
	
	// **[RUTE CRUD UMUM - Diletakkan Setelah Rute Spesifik/Berparameter]**

	// CREATE Pekerjaan baru (POST /unair/pekerjaan)
	protected.Post("/", jobService.Create)
	
	// UPDATE Pekerjaan milik sendiri (PUT /unair/pekerjaan/:id)
	protected.Put("/:id", jobService.Update)

	// SOFT DELETE Pekerjaan milik sendiri (DELETE /unair/pekerjaan/:id)
	protected.Delete("/:id", jobService.SoftDeleteByAlumni) 

	// HARD DELETE Pekerjaan milik sendiri (DELETE /unair/pekerjaan/hard/:id)
	protected.Delete("/hard/:id", jobService.HardDelete)

	// === [1. PUBLIC ROUTES - Akses Tanpa Autentikasi] ===
	// Rute untuk melihat data pekerjaan yang aktif secara publik
	job.Get("/", jobService.GetAll)
	job.Get("/:id", jobService.GetByID) // <- Ditempatkan paling akhir di group `job`
}
