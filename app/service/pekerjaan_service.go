// package service

// import (
// 	"clean-archi/app/model"
// 	"clean-archi/app/repository"
// 	"database/sql"
// 	"fmt"
// 	"strconv"

// 	"github.com/gofiber/fiber/v2"
// )

// // JobService menyimpan dependency database untuk digunakan di semua method
// type JobService struct {
// 	db *sql.DB
// }

// // NewJobService membuat instance baru JobService dengan db yang diberikan
// func NewJobService(db *sql.DB) *JobService {
// 	return &JobService{db: db}
// }

// // Helper function: mengambil userID dari Locals (JWT token)
// // Helper function untuk mengambil userID dari Locals dengan aman
// func getUserIDFromLocals(c *fiber.Ctx) (int, error) {
// 	userIDValue := c.Locals("user_id")
// 	if userIDValue == nil {
// 		return 0, fmt.Errorf("user ID tidak ditemukan di konteks")
// 	}

// 	// Cek jika ID adalah int (tipe data native Go/DB)
// 	if idInt, ok := userIDValue.(int); ok {
// 		return idInt, nil
// 	}
// 	// Cek jika ID adalah float64 (umum terjadi saat decode dari JWT/JSON)
// 	if idFloat, ok := userIDValue.(float64); ok {
// 		return int(idFloat), nil
// 	}

// 	return 0, fmt.Errorf("tipe data user ID tidak valid")
// }

// // ========================
// // GET METHODS
// // ========================

// // GetAll: ambil semua pekerjaan aktif (is_deleted = NULL)
// func (s *JobService) GetAll(c *fiber.Ctx) error {
// 	// ... (kode GetAll tidak berubah) ...
// 	jobs, err := repository.GetAllJobs(s.db)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
// 			"success": false,
// 		})
// 	}

// 	if len(jobs) == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"message": "Tidak ada data pekerjaan yang ditemukan",
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Berhasil mendapatkan semua data pekerjaan",
// 		"success": true,
// 		"data": 	 jobs,
// 	})
// }

// // GetByID: ambil pekerjaan berdasarkan ID, cek soft delete
// func (s *JobService) GetByID(c *fiber.Ctx) error {
// 	// ... (kode GetByID tidak berubah) ...
// 	id, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID tidak valid", "success": false})
// 	}

// 	job, err := repository.GetJobByID(s.db, id)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error: " + err.Error(), "success": false})
// 	}
// 	if job == nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Pekerjaan tidak ditemukan", "success": false})
// 	}
// 	if job.IsDeleted != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Data pekerjaan telah dihapus", "success": false})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Berhasil mendapatkan data pekerjaan", "success": true, "data": job})
// }

// // GetJobsByAlumniID: ambil semua pekerjaan milik alumni tertentu
// func (s *JobService) GetJobsByAlumniID(c *fiber.Ctx) error {
// 	// ... (kode GetJobsByAlumniID tidak berubah) ...
// 	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "alumni_id tidak valid", "success": false})
// 	}

// 	jobs, err := repository.GetJobsByAlumniID(s.db, alumniID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data pekerjaan: " + err.Error(), "success": false})
// 	}
// 	if len(jobs) == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Tidak ada pekerjaan untuk alumni ini", "success": false})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Berhasil mendapatkan data pekerjaan", "success": true, "data": jobs})
// }

// // GetMyJobs: ambil pekerjaan milik user yang login (ID dari JWT)
// // GetMyJobs menangani permintaan GET /unair/pekerjaan/my. Mengambil semua pekerjaan aktif milik user yang sedang login.
// // Handler ini membaca ID dari Locals (token JWT). Ini adalah versi clean dari GetJobsByAlumniID untuk rute /my.
// func (s *JobService) GetMyJobs(c *fiber.Ctx) error {
// 	// Ambil ID dari Locals (token JWT)
// 	alumniID, err := getUserIDFromLocals(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Akses ditolak: User ID tidak ditemukan atau tidak valid", "success": false})
// 	}

// 	// Panggil repository dengan ID yang sudah pasti aman dari token
// 	jobs, err := repository.GetJobsByAlumniID(s.db, alumniID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data pekerjaan: " + err.Error(), "success": false})
// 	}
// 	if len(jobs) == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Anda belum memiliki data pekerjaan", "success": false})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Berhasil mendapatkan data pekerjaan milik Anda", "success": true, "data": jobs})
// }

// // ========================
// // CREATE & UPDATE
// // ========================

// // Create: tambah pekerjaan baru
// func (s *JobService) Create(c *fiber.Ctx) error {
// 	// ... (kode Create tidak berubah) ...
// 	var pekerjaan model.Pekerjaan
// 	if err := c.BodyParser(&pekerjaan); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body tidak valid", "success": false})
// 	}

// 	created, err := repository.CreateJob(s.db, &pekerjaan)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menambahkan pekerjaan: " + err.Error(), "success": false})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Pekerjaan berhasil ditambahkan", "success": true, "data": created})
// }

// // Update: update pekerjaan
// func (s *JobService) Update(c *fiber.Ctx) error {
// 	// ... (kode Update tidak berubah) ...
// 	id, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID tidak valid", "success": false})
// 	}

// 	var pekerjaan model.Pekerjaan
// 	if err := c.BodyParser(&pekerjaan); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request body tidak valid", "success": false})
// 	}

// 	updated, err := repository.UpdateJob(s.db, id, &pekerjaan)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengupdate pekerjaan: " + err.Error(), "success": false})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Pekerjaan berhasil diperbarui", "success": true, "data": updated})
// }

// // ========================
// // SOFT DELETE
// // ========================

// // SoftDeleteByAlumni: alumni hanya bisa hapus pekerjaannya sendiri
// func (s *JobService) SoftDeleteByAlumni(c *fiber.Ctx) error {
// 	// FIX: Menggunakan helper function
// 	userID, err := getUserIDFromLocals(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error(), "success": false})
// 	}

// 	jobID, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "ID pekerjaan tidak valid",
// 			"success": false,
// 		})
// 	}

// 	err = repository.SoftDeletePekerjaanByAlumni(s.db, jobID, userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": err.Error(),
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Pekerjaan berhasil dihapus (soft delete)",
// 		"success": true,
// 	})
// }

// // SoftDeleteAllByAdmin: admin bisa hapus semua pekerjaan alumni tertentu
// func (s *JobService) SoftDeleteAllByAdmin(c *fiber.Ctx) error {
// 	role := c.Locals("role").(string)
// 	if role != "admin" {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
// 			"message": "Hanya admin yang dapat melakukan ini",
// 			"success": false,
// 		})
// 	}

// 	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "ID alumni tidak valid",
// 			"success": false,
// 		})
// 	}

// 	count, err := repository.SoftDeleteAllPekerjaanByAdmin(s.db, alumniID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": err.Error(),
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": 		 fmt.Sprintf("Berhasil menghapus %d pekerjaan milik alumni ID %d", count, alumniID),
// 		"deleted_count": count,
// 		"success": 		 true,
// 	})
// }

// // ========================
// // TRASHED & RESTORE
// // ========================

// // GetTrashed: ambil pekerjaan yang dihapus
// func (s *JobService) GetTrashed(c *fiber.Ctx) error {
// 	role := c.Locals("role").(string)

// 	// FIX: Menggunakan helper function
// 	userID, err := getUserIDFromLocals(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error(), "success": false})
// 	}

// 	var jobs []model.Pekerjaan

// 	if role == "admin" {
// 		jobs, err = repository.GetTrashedJobs(s.db)
// 	} else {
// 		jobs, err = repository.GetTrashedJobsByAlumni(s.db, userID)
// 	}

// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": err.Error(),
// 			"success": false,
// 		})
// 	}

// 	if len(jobs) == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"message": "Tidak ada data yang dihapus",
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Data pekerjaan yang dihapus berhasil diambil",
// 		"data": 	 jobs,
// 		"success": true,
// 	})
// }

// // ========================
// // RESTORE & HARD DELETE
// // ========================

// // Restore: kembalikan pekerjaan yang di-soft delete
// // Admin bisa restore semua pekerjaan, alumni hanya miliknya
// func (s *JobService) Restore(c *fiber.Ctx) error {
// 	role := c.Locals("role").(string)

// 	// FIX: Menggunakan helper function
// 	userID, err := getUserIDFromLocals(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error(), "success": false})
// 	}

// 	jobID, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "ID pekerjaan tidak valid",
// 			"success": false,
// 		})
// 	}

// 	if role == "admin" {
// 		err = repository.RestoreJob(s.db, jobID)
// 	} else {
// 		err = repository.RestoreJobByAlumni(s.db, jobID, userID)
// 	}

// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": err.Error(),
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Pekerjaan berhasil direstore",
// 		"success": true,
// 	})
// }

// // HardDelete: hapus permanen pekerjaan
// // Admin bisa hapus semua pekerjaan, alumni hanya miliknya
// func (s *JobService) HardDelete(c *fiber.Ctx) error {
// 	//untuk mengambil peran user dari context request
// 	// supaya handler tahu apa yang boleh dilakukan user
// 	role := c.Locals("role").(string)

// 	// FIX: Menggunakan helper function
// 	userID, err := getUserIDFromLocals(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error(), "success": false})
// 	}

// 	jobID, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "ID pekerjaan tidak valid",
// 			"success": false,
// 		})
// 	}

// 	if role == "admin" {
// 		err = repository.HardDeleteJob(s.db, jobID)
// 	} else {
// 		err = repository.HardDeleteJobByAlumni(s.db, jobID, userID)
// 	}

// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": err.Error(),
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Pekerjaan berhasil dihapus permanen",
// 		"success": true,
// 	})
// }

package service

import (
	"clean-archi/app/model"
	"clean-archi/app/repository"
	"fmt"
	"time"


	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============================================================
// ================        JOB SERVICE         =================
// ============================================================

// JobService menangani operasi CRUD untuk data pekerjaan menggunakan MongoDB.
type JobService struct {
    repo repository.PekerjaanRepository
}


// NewJobService membuat instance baru dari JobService.
func NewJobService(repo repository.PekerjaanRepository) *JobService {
    return &JobService{repo: repo}
}


// ============================================================
// ================         HELPER METHOD      =================
// ============================================================

// getUserIDFromLocals mengambil user ID (string ObjectID) dari context.
func getUserIDFromLocals(c *fiber.Ctx) (string, error) {
	userIDValue := c.Locals("user_id") // Asumsi Auth middleware menyimpan string ObjectID
	if userIDValue == nil {
		return "", fmt.Errorf("user ID tidak ditemukan di konteks")
	}

	userID, ok := userIDValue.(string)
	if !ok {
		// Jika tersimpan sebagai primitive.ObjectID, konversi ke string
		if oid, oidOK := userIDValue.(primitive.ObjectID); oidOK {
			return oid.Hex(), nil
		}
		return "", fmt.Errorf("tipe data user ID tidak valid (bukan string ObjectID)")
	}
	return userID, nil
}

// ============================================================
// ================         CRUD HANDLERS      =================
// ============================================================

// GetAll godoc
// @Summary Ambil semua data pekerjaan
// @Description Mengambil seluruh daftar pekerjaan dari MongoDB
// @Tags Pekerjaan
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Pekerjaan
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /pekerjaan [get]
func (s *JobService) GetAll(c *fiber.Ctx) error {
	jobs, err := s.repo.GetAll()
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
		"message": "Berhasil mendapatkan semua pekerjaan",
		"success": true,
		"data":    jobs,
	})
}

// @Security BearerAuth
// GetByID godoc
// @Summary Ambil data pekerjaan berdasarkan ID
// @Description Mengambil data pekerjaan dari MongoDB berdasarkan ID yang diberikan
// @Tags Pekerjaan
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Success 200 {object} model.Pekerjaan
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /pekerjaan/{id} [get]
func (s *JobService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID diperlukan",
			"success": false,
		})
	}

	job, err := s.repo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil pekerjaan: " + err.Error(),
			"success": false,
		})
	}
	if job == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pekerjaan tidak ditemukan",
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data pekerjaan",
		"success": true,
		"data":    job,
	})
}

// Create godoc
// @Security BearerAuth
// @Summary Tambah pekerjaan baru
// @Description Menambahkan data pekerjaan baru ke MongoDB
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param pekerjaan body model.CreatePekerjaanRequest true "Data pekerjaan baru"
// @Success 201 {object} model.Pekerjaan
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /pekerjaan [post]
func (s *JobService) Create(c *fiber.Ctx) error {
	var req model.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body tidak valid",
			"success": false,
		})
	}

	// 🧩 Konversi array string AlumniID dari request ke array ObjectID
	var alumniObjIDs []primitive.ObjectID
	for _, idStr := range req.AlumniIDs {
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "alumni_id tidak valid: " + idStr,
				"success": false,
			})
		}
		alumniObjIDs = append(alumniObjIDs, oid)
	}

	// 🧩 Buat data pekerjaan baru
	pekerjaan := model.Pekerjaan{
		AlumniIDs:          alumniObjIDs,
		NamaPerusahaan:     req.NamaPerusahaan,
		PosisiJabatan:      req.PosisiJabatan,
		BidangIndustri:     req.BidangIndustri,
		LokasiKerja:        req.LokasiKerja,
		DeskripsiPekerjaan: req.DeskripsiPekerjaan,
		TanggalMulaiKerja:  req.TanggalMulaiKerja,
		TanggalSelesaiKerja: req.TanggalSelesaiKerja,
		StatusPekerjaan:    req.StatusPekerjaan,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		IsDeleted:          false,
	}

	// 🧩 Simpan ke MongoDB
	if err := s.repo.Create(&pekerjaan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	// 🧩 Respons sukses
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pekerjaan berhasil ditambahkan",
		"success": true,
		"data":    pekerjaan,
	})
}

// Update godoc
// @Security BearerAuth
// @Summary Perbarui data pekerjaan
// @Description Memperbarui data pekerjaan yang ada di MongoDB berdasarkan ID
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Param pekerjaan body model.Pekerjaan true "Data pekerjaan yang diperbarui"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /pekerjaan/{id} [put]
func (s *JobService) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID diperlukan",
			"success": false,
		})
	}

	var pekerjaan model.Pekerjaan
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request body tidak valid",
			"success": false,
		})
	}

	pekerjaan.UpdatedAt = time.Now()

	if err := s.repo.Update(id, &pekerjaan); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal memperbarui pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil diperbarui",
		"success": true,
	})
}

// Delete godoc
// @Security BearerAuth
// @Summary Hapus data pekerjaan
// @Description Menghapus data pekerjaan dari MongoDB berdasarkan ID
// @Tags Pekerjaan
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /pekerjaan/{id} [delete]
func (s *JobService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID diperlukan",
			"success": false,
		})
	}

	if err := s.repo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menghapus pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan berhasil dihapus",
		"success": true,
	})
}

// ============================================================
// SOFT DELETE, TRASH, RESTORE, HARD DELETE
// ============================================================

// SoftDeleteByAlumni godoc
// @Security BearerAuth
// @Summary Hapus pekerjaan (Soft Delete) oleh Alumni
// @Description Alumni dapat menghapus pekerjaan mereka sendiri, data tidak benar-benar dihapus dari database (soft delete).
// @Tags Pekerjaan
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /pekerjaan/{id}/soft-delete [delete]
func (s *JobService) SoftDeleteByAlumni(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	jobID := c.Params("id")
	if jobID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "ID pekerjaan tidak valid", "success": false})
	}

	if err := s.repo.SoftDeleteByAlumni(jobID, userID); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus (soft delete)", "success": true})
}

// SoftDeleteAllByAdmin godoc
// @Security BearerAuth
// @Summary Hapus semua pekerjaan milik alumni (Soft Delete) oleh Admin
// @Description Admin dapat melakukan soft delete terhadap semua pekerjaan milik satu alumni berdasarkan ID alumni.
// @Tags Pekerjaan
// @Produce json
// @Param alumni_id path string true "ID Alumni (ObjectID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /pekerjaan/soft-delete/all/{alumni_id} [delete]
func (s *JobService) SoftDeleteAllByAdmin(c *fiber.Ctx) error {
	role, ok := c.Locals("role").(string)
	if !ok || role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "Hanya admin yang dapat melakukan ini",
		})
	}

	adminID, err := getUserIDFromLocals(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Gagal membaca admin ID", "success": false})
	}

	alumniID := c.Params("alumni_id")
	if alumniID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "alumni_id wajib diisi",
		})
	}

	count, err := s.repo.SoftDeleteAllByAdmin(alumniID, adminID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Gagal melakukan soft delete: %v", err),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Berhasil menghapus %d pekerjaan milik alumni ID %s", count, alumniID),
	})
}

// GetTrased godoc
// @Security BearerAuth
// @Summary Ambil daftar pekerjaan yang telah dihapus (soft deleted)
// @Description Mengambil daftar pekerjaan yang telah di-soft delete. Admin melihat semua data, Alumni hanya melihat data miliknya.
// @Tags Pekerjaan
// @Produce json
// @Success 200 {object} map[string]interface{} "Berhasil mengambil data yang dihapus"
// @Failure 404 {object} map[string]interface{} "Tidak ada data yang dihapus"
// @Failure 500 {object} map[string]interface{} "Kesalahan Internal Server"
// @Router /pekerjaan/trashed [get]
func (s *JobService) GetTrashed(c *fiber.Ctx) error {
	role, _ := c.Locals("role").(string)
	userID, _ := getUserIDFromLocals(c)

	// 🧩 Tambahkan di sini
	fmt.Println("DEBUG userID:", userID)

	var jobs []model.Pekerjaan
	var err error

	if role == "admin" {
		// ✅ Admin: ambil semua pekerjaan yang dihapus
		jobs, err = s.repo.GetTrashedJobs()
	} else {
		// ✅ User biasa: validasi ObjectID dulu
		if !primitive.IsValidObjectID(userID) {
			return c.Status(400).JSON(fiber.Map{
				"message": "User ID tidak valid (bukan ObjectID)",
				"success": false,
			})
		}
		jobs, err = s.repo.GetTrashedJobsByAlumni(userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil pekerjaan: " + err.Error(),
			"success": false,
		})
	}

	if len(jobs) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Tidak ada data yang dihapus",
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil data yang dihapus",
		"success": true,
		"data":    jobs,
	})
}


// Restore godoc
// @Security BearerAuth
// @Summary Restore pekerjaan (mengembalikan status aktif)
// @Description Mengembalikan status pekerjaan dari soft deleted menjadi aktif. Admin dapat restore data siapa pun, Alumni hanya data miliknya.
// @Tags Pekerjaan
// @Produce json
// @Param id path string true "ID Pekerjaan (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Pekerjaan berhasil direstore"
// @Failure 400 {object} map[string]interface{} "Gagal restore (misal: data tidak ditemukan atau bukan milik pengguna)"
// @Router /pekerjaan/restore/{id} [put]
func (s *JobService) Restore(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID, _ := getUserIDFromLocals(c)
	jobID := c.Params("id")

	var err error
	if role == "admin" {
		err = s.repo.RestoreJob(jobID)
	} else {
		err = s.repo.RestoreJobByAlumni(jobID, userID)
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error(), "success": false})
	}
	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil direstore", "success": true})
}

// HardDelete godoc
// @Security BearerAuth
// @Summary Hapus permanen pekerjaan (Hard Delete)
// @Description Menghapus data pekerjaan secara permanen dari database. Admin dapat menghapus data siapa pun, Alumni hanya data miliknya.
// @Tags Pekerjaan
// @Produce json
// @Param id path string true "ID Pekerjaan (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{} "Pekerjaan berhasil dihapus permanen"
// @Failure 400 {object} map[string]interface{} "Gagal hapus permanen (misal: data tidak ditemukan atau bukan milik pengguna)"
// @Router /pekerjaan/hard-delete/{id} [delete]
func (s *JobService) HardDelete(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID, _ := getUserIDFromLocals(c)
	jobID := c.Params("id")

	var err error
	if role == "admin" {
		err = s.repo.HardDeleteJob(jobID)
	} else {
		err = s.repo.HardDeleteJobByAlumni(jobID, userID)
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error(), "success": false})
	}
	return c.JSON(fiber.Map{"message": "Pekerjaan berhasil dihapus permanen", "success": true})
}
