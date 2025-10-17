package service

import (
	"clean-archi/app/model"
	"clean-archi/app/repository"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// JobService menyimpan dependency database untuk digunakan di semua method
type JobService struct {
	db *sql.DB
}

// NewJobService membuat instance baru JobService dengan db yang diberikan
func NewJobService(db *sql.DB) *JobService {
	return &JobService{db: db}
}

// Helper function: mengambil userID dari Locals (JWT token)
// Helper function untuk mengambil userID dari Locals dengan aman
func getUserIDFromLocals(c *fiber.Ctx) (int, error) {
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return 0, fmt.Errorf("user ID tidak ditemukan di konteks")
	}

	// Cek jika ID adalah int (tipe data native Go/DB)
	if idInt, ok := userIDValue.(int); ok {
		return idInt, nil
	}
	// Cek jika ID adalah float64 (umum terjadi saat decode dari JWT/JSON)
	if idFloat, ok := userIDValue.(float64); ok {
		return int(idFloat), nil
	}
	
	return 0, fmt.Errorf("tipe data user ID tidak valid")
}

// ========================
// GET METHODS
// ========================

// GetAll: ambil semua pekerjaan aktif (is_deleted = NULL)
func (s *JobService) GetAll(c *fiber.Ctx) error {
	// ... (kode GetAll tidak berubah) ...
	jobs, err := repository.GetAllJobs(s.db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	if len(jobs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Tidak ada data pekerjaan yang ditemukan",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mendapatkan semua data pekerjaan",
		"success": true,
		"data": 	 jobs,
	})
}


// GetByID: ambil pekerjaan berdasarkan ID, cek soft delete
func (s *JobService) GetByID(c *fiber.Ctx) error {
	// ... (kode GetByID tidak berubah) ...
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID tidak valid", "success": false})
	}

	job, err := repository.GetJobByID(s.db, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error: " + err.Error(), "success": false})
	}
	if job == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pekerjaan tidak ditemukan", "success": false})
	}
	if job.IsDeleted != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data pekerjaan telah dihapus", "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Berhasil mendapatkan data pekerjaan", "success": true, "data": job})
}

// GetJobsByAlumniID: ambil semua pekerjaan milik alumni tertentu
func (s *JobService) GetJobsByAlumniID(c *fiber.Ctx) error {
	// ... (kode GetJobsByAlumniID tidak berubah) ...
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "alumni_id tidak valid", "success": false})
	}

	jobs, err := repository.GetJobsByAlumniID(s.db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data pekerjaan: " + err.Error(), "success": false})
	}
	if len(jobs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Tidak ada pekerjaan untuk alumni ini", "success": false})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Berhasil mendapatkan data pekerjaan", "success": true, "data": jobs})
}

// GetMyJobs: ambil pekerjaan milik user yang login (ID dari JWT)
// GetMyJobs menangani permintaan GET /unair/pekerjaan/my. Mengambil semua pekerjaan aktif milik user yang sedang login.
// Handler ini membaca ID dari Locals (token JWT). Ini adalah versi clean dari GetJobsByAlumniID untuk rute /my.
func (s *JobService) GetMyJobs(c *fiber.Ctx) error {
	// Ambil ID dari Locals (token JWT)
	alumniID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Akses ditolak: User ID tidak ditemukan atau tidak valid", "success": false})
	}
	
	// Panggil repository dengan ID yang sudah pasti aman dari token
	jobs, err := repository.GetJobsByAlumniID(s.db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data pekerjaan: " + err.Error(), "success": false})
	}
	if len(jobs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Anda belum memiliki data pekerjaan", "success": false})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Berhasil mendapatkan data pekerjaan milik Anda", "success": true, "data": jobs})
}

// ========================
// CREATE & UPDATE
// ========================


// Create: tambah pekerjaan baru
func (s *JobService) Create(c *fiber.Ctx) error {
	// ... (kode Create tidak berubah) ...
	var pekerjaan model.Pekerjaan
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body tidak valid", "success": false})
	}

	created, err := repository.CreateJob(s.db, &pekerjaan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menambahkan pekerjaan: " + err.Error(), "success": false})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Pekerjaan berhasil ditambahkan", "success": true, "data": created})
}

// Update: update pekerjaan
func (s *JobService) Update(c *fiber.Ctx) error {
	// ... (kode Update tidak berubah) ...
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID tidak valid", "success": false})
	}

	var pekerjaan model.Pekerjaan
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body tidak valid", "success": false})
	}

	updated, err := repository.UpdateJob(s.db, id, &pekerjaan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengupdate pekerjaan: " + err.Error(), "success": false})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Pekerjaan berhasil diperbarui", "success": true, "data": updated})
}

// ========================
// SOFT DELETE
// ========================


// SoftDeleteByAlumni: alumni hanya bisa hapus pekerjaannya sendiri
func (s *JobService) SoftDeleteByAlumni(c *fiber.Ctx) error {
	// FIX: Menggunakan helper function
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	jobID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID pekerjaan tidak valid",
			"success": false,
		})
	}

	err = repository.SoftDeletePekerjaanByAlumni(s.db, jobID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil dihapus (soft delete)",
		"success": true,
	})
}

// SoftDeleteAllByAdmin: admin bisa hapus semua pekerjaan alumni tertentu
func (s *JobService) SoftDeleteAllByAdmin(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Hanya admin yang dapat melakukan ini",
			"success": false,
		})
	}

	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID alumni tidak valid",
			"success": false,
		})
	}

	count, err := repository.SoftDeleteAllPekerjaanByAdmin(s.db, alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": 		 fmt.Sprintf("Berhasil menghapus %d pekerjaan milik alumni ID %d", count, alumniID),
		"deleted_count": count,
		"success": 		 true,
	})
}

// ========================
// TRASHED & RESTORE
// ========================


// GetTrashed: ambil pekerjaan yang dihapus
func (s *JobService) GetTrashed(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	
	// FIX: Menggunakan helper function
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	var jobs []model.Pekerjaan
	
	if role == "admin" {
		jobs, err = repository.GetTrashedJobs(s.db)
	} else {
		jobs, err = repository.GetTrashedJobsByAlumni(s.db, userID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"success": false,
		})
	}

	if len(jobs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Tidak ada data yang dihapus",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data pekerjaan yang dihapus berhasil diambil",
		"data": 	 jobs,
		"success": true,
	})
}

// ========================
// RESTORE & HARD DELETE
// ========================

// Restore: kembalikan pekerjaan yang di-soft delete
// Admin bisa restore semua pekerjaan, alumni hanya miliknya
func (s *JobService) Restore(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	
	// FIX: Menggunakan helper function
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error(), "success": false})
	}
	
	jobID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID pekerjaan tidak valid",
			"success": false,
		})
	}

	if role == "admin" {
		err = repository.RestoreJob(s.db, jobID)
	} else {
		err = repository.RestoreJobByAlumni(s.db, jobID, userID)
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil direstore",
		"success": true,
	})
}

// HardDelete: hapus permanen pekerjaan
// Admin bisa hapus semua pekerjaan, alumni hanya miliknya
func (s *JobService) HardDelete(c *fiber.Ctx) error {
	//untuk mengambil peran user dari context request 
	// supaya handler tahu apa yang boleh dilakukan user
	role := c.Locals("role").(string)
	
	// FIX: Menggunakan helper function
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error(), "success": false})
	}
	
	jobID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID pekerjaan tidak valid",
			"success": false,
		})
	}

	if role == "admin" {
		err = repository.HardDeleteJob(s.db, jobID)
	} else {
		err = repository.HardDeleteJobByAlumni(s.db, jobID, userID)
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil dihapus permanen",
		"success": true,
	})
}
