package routes

import (
    "clean-archi/app/service"
    repo "clean-archi/app/repository/MongoRepo"
    "clean-archi/middleware"
    "os"

    "go.mongodb.org/mongo-driver/mongo"
    "github.com/gofiber/fiber/v2"
)

// JobRoutes sekarang hanya menggunakan MongoDB
func JobRoutes(router fiber.Router, mongoClient *mongo.Client) {
    // Ambil database dari nama di .env
    mongoDatabase := mongoClient.Database(os.Getenv("MONGODB_DATABASE"))

    // Inisialisasi repository untuk koleksi "pekerjaan"
    mongoRepo := repo.NewPekerjaanMongoRepository(mongoDatabase, "pekerjaan")

    // Inisialisasi service
    jobService := service.NewJobService(mongoRepo)

    // Base group: /unair/pekerjaan
    job := router.Group("/api/v1/pekerjaan")

    // Middleware auth (jika masih dipakai)
    protected := job.Group("", middleware.AuthRequired())

    // ==== PUBLIC ROUTES ====
    job.Get("/", jobService.GetAll)
    job.Get("/:id", jobService.GetByID)

    // ==== PROTECTED ROUTES ====
    protected.Post("/", jobService.Create)
    protected.Put("/:id", jobService.Update)
    protected.Delete("/:id", jobService.Delete)

    // ============================================================
	// ==== TAMBAHAN RUTE (SOFT DELETE, TRASH, RESTORE) ====
	// ============================================================

	// Rute untuk Soft Delete oleh Alumni
	// DELETE /api/v1/pekerjaan/{id}/soft-delete
	protected.Delete("/:id/soft-delete", jobService.SoftDeleteByAlumni)

	// Rute untuk Soft Delete Semua Pekerjaan oleh Admin (berdasarkan ID Alumni)
	// DELETE /api/v1/pekerjaan/soft-delete/all/{alumni_id}
	protected.Delete("/soft-delete/all/:alumni_id", jobService.SoftDeleteAllByAdmin)

	// Rute untuk melihat data yang di-trash
	// GET /api/v1/pekerjaan/trashed
	protected.Get("/trashed", jobService.GetTrashed)

	// Rute untuk Restore data
	// PUT /api/v1/pekerjaan/restore/{id}
	protected.Put("/restore/:id", jobService.Restore)

	// Rute untuk Hard Delete (sebaiknya dibedakan dari CRUD Delete biasa)
	// DELETE /api/v1/pekerjaan/hard-delete/{id}
	protected.Delete("/hard-delete/:id", jobService.HardDelete)
}

