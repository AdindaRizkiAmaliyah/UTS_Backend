package service

import (
	"clean-archi/app/model"
	"clean-archi/app/repository/MongoRepo"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ============================================================
// ================         STRUCT SERVICE     =================
// ============================================================

// FileService menangani proses upload dan manajemen file (foto, sertifikat, dsb)
type FileService struct {
	repo       *MongoRepo.FileRepository
	uploadPath string
}

// NewFileService membuat instance baru dari FileService
func NewFileService(repo *MongoRepo.FileRepository, uploadPath string) *FileService {
	return &FileService{repo: repo, uploadPath: uploadPath}
}

// ============================================================
// ================         UPLOAD FOTO        =================
// ============================================================

// UploadFoto godoc
// @Security BearerAuth
// @Summary Upload foto pengguna
// @Description Hanya admin atau pemilik akun yang dapat mengunggah foto profil. Format file: JPG, JPEG, atau PNG (maks 1MB).
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData string true "ID pengguna"
// @Param file formData file true "File foto"
// @Success 201 {object} model.File
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /files/upload/foto [post]
func (s *FileService) UploadFoto(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid atau user tidak ditemukan",
		})
	}

	userID := c.FormValue("user_id")
	if user.Role != "admin" && userID != user.UserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Hanya admin atau pemilik akun yang bisa upload",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File tidak ditemukan")
	}

	if fileHeader.Size > 1*1024*1024 {
		return fiber.NewError(fiber.StatusBadRequest, "Ukuran file maksimal 1MB")
	}

	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}
	if !allowed[fileHeader.Header.Get("Content-Type")] {
		return fiber.NewError(fiber.StatusBadRequest, "Format file foto tidak diizinkan")
	}

	return s.saveFile(c, fileHeader, userID, "foto")
}

// ============================================================
// ================      UPLOAD SERTIFIKAT     =================
// ============================================================

// UploadSertifikat godoc
// @Security BearerAuth
// @Summary Upload sertifikat (PDF)
// @Description Hanya admin atau pemilik akun yang dapat mengunggah file sertifikat dalam format PDF (maks 2MB).
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData string true "ID pengguna"
// @Param file formData file true "File sertifikat (PDF)"
// @Success 201 {object} model.File
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /files/upload/sertifikat [post]
func (s *FileService) UploadSertifikat(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid atau user tidak ditemukan",
		})
	}

	userID := c.FormValue("user_id")
	if user.Role != "admin" && userID != user.UserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Hanya admin atau pemilik akun yang bisa upload",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File tidak ditemukan")
	}

	if fileHeader.Size > 2*1024*1024 {
		return fiber.NewError(fiber.StatusBadRequest, "Ukuran file maksimal 2MB")
	}

	if fileHeader.Header.Get("Content-Type") != "application/pdf" {
		return fiber.NewError(fiber.StatusBadRequest, "Hanya file PDF yang diperbolehkan")
	}

	return s.saveFile(c, fileHeader, userID, "sertifikat")
}

// ============================================================
// ================       SAVE FILE HELPER     =================
// ============================================================

// saveFile menyimpan file ke folder upload dan mencatat metadata-nya ke database
func (s *FileService) saveFile(c *fiber.Ctx, fileHeader *multipart.FileHeader, userID, category string) error {
	ext := filepath.Ext(fileHeader.Filename)
	newName := uuid.New().String() + ext
	dir := filepath.Join(s.uploadPath, category)

	// Pastikan folder upload ada
	os.MkdirAll(dir, os.ModePerm)
	path := filepath.Join(dir, newName)

	// Simpan file fisik
	if err := c.SaveFile(fileHeader, path); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan file")
	}

	// Simpan data file ke database
	fileModel := &model.File{
		UserID:       userID,
		FileName:     newName,
		OriginalName: fileHeader.Filename,
		FilePath:     path,
		FileSize:     fileHeader.Size,
		FileType:     fileHeader.Header.Get("Content-Type"),
		FileCategory: category,
	}

	if err := s.repo.Create(fileModel); err != nil {
		os.Remove(path)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan metadata file")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("%s berhasil diunggah", category),
		"data":    fileModel,
	})
}

// ============================================================
// ================         CRUD TAMBAHAN     =================
// ============================================================

// GetAllFiles godoc
// @Security BearerAuth
// @Summary Ambil semua file (admin only)
// @Description Mengambil seluruh data file yang tersimpan dalam sistem
// @Tags File
// @Produce json
// @Success 200 {array} model.File
// @Failure 500 {object} map[string]interface{}
// @Router /files/ [get]
func (s *FileService) GetAllFiles(c *fiber.Ctx) error {
	files, err := s.repo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil semua file",
		})
	}
	return c.JSON(files)
}

// GetFileByID godoc
// @Security BearerAuth
// @Summary Ambil detail file berdasarkan ID
// @Description Mendapatkan data file berdasarkan ID yang diberikan
// @Tags File
// @Produce json
// @Param id path string true "ID file"
// @Success 200 {object} model.File
// @Failure 404 {object} map[string]interface{}
// @Router /api/files/{id} [get]
func (s *FileService) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "File tidak ditemukan",
		})
	}
	return c.JSON(file)
}

// DeleteFile godoc
// @Security BearerAuth
// @Summary Hapus file berdasarkan ID
// @Description Hanya admin yang dapat menghapus file dari sistem
// @Tags File
// @Produce json
// @Param id path string true "ID file"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/files/{id} [delete]
func (s *FileService) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.repo.DeleteByID(id); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "File tidak ditemukan",
		})
	}
	return c.JSON(fiber.Map{
		"message": "File berhasil dihapus",
	})
}

// GetMyFiles godoc
// @Security BearerAuth
// @Summary Ambil semua file milik user yang sedang login
// @Description Menampilkan daftar file milik user yang sesuai dengan token login aktif
// @Tags File
// @Produce json
// @Success 200 {array} model.File
// @Failure 500 {object} map[string]interface{}
// @Router /file/me [get]
func (s *FileService) GetMyFiles(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid atau user tidak ditemukan",
		})
	}

	files, err := s.repo.GetByUserID(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil file milik user",
		})
	}
	return c.JSON(files)
}
