package routes

import (
	"clean-archi/middleware"
	"clean-archi/app/service"
	"github.com/gofiber/fiber/v2"
)


// FileRoutes — semua endpoint file (foto, sertifikat, dan manajemen file)
func FileRoutes(app *fiber.App, fileService *service.FileService) {
	// Semua endpoint di bawah wajib login
	api := app.Group("/api/v1/files", middleware.AuthRequired())

	// Upload Foto
	api.Post("/upload/foto", fileService.UploadFoto)

	// Upload Sertifikat
	api.Post("/upload/sertifikat", fileService.UploadSertifikat)

	// Admin-only section
	admin := api.Group("/", middleware.AdminOnly())
	admin.Get("/", fileService.GetAllFiles)
	admin.Get("/:id", fileService.GetFileByID)
	admin.Delete("/:id", fileService.DeleteFile)

	// User melihat file miliknya
	api.Get("/my", fileService.GetMyFiles)
}
