package service

import (
	"clean-archi/app/model"
	"clean-archi/app/repository"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"fmt"
)

type JobService struct {
	db *sql.DB
}

func NewJobService(db *sql.DB) *JobService {
	return &JobService{db: db}
}

// ✅ GetAll — dengan dukungan filter query params
func (s *JobService) GetAll(c *fiber.Ctx) error {
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
		"data":    jobs,
	})
}


// ✅ GetByID
func (s *JobService) GetByID(c *fiber.Ctx) error {
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

// ✅ GetJobsByAlumniID
func (s *JobService) GetJobsByAlumniID(c *fiber.Ctx) error {
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

// ✅ Create
func (s *JobService) Create(c *fiber.Ctx) error {
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

// ✅ Update
func (s *JobService) Update(c *fiber.Ctx) error {
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

// 🧱 Soft delete pekerjaan oleh alumni
func (s *JobService) SoftDeleteByAlumni(c *fiber.Ctx) error {
    // Ambil user_id dari token
    userID, ok := c.Locals("user_id").(int)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "User ID tidak valid dalam token",
            "success": false,
        })
    }

    // Ambil ID pekerjaan dari parameter
    jobID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "ID pekerjaan tidak valid",
            "success": false,
        })
    }

    // Lanjutkan soft delete di database
    query := `UPDATE pekerjaan SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND id_alumni = ?`
    result, err := s.DB.Exec(query, jobID, userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menghapus pekerjaan",
            "success": false,
        })
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "Pekerjaan tidak ditemukan atau bukan milik user ini",
            "success": false,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Pekerjaan berhasil dihapus (soft delete)",
        "success": true,
    })
}


// 🧱 Soft delete semua pekerjaan alumni oleh admin
func (s *JobService) SoftDeleteAllByAdmin(c *fiber.Ctx) error {
	alumniID, _ := strconv.Atoi(c.Params("alumni_id"))
	adminID := int(c.Locals("user_id").(float64))
	role := c.Locals("role").(string)

	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"message": "Hanya admin yang dapat melakukan ini", "success": false})
	}

	count, err := repository.SoftDeleteAllPekerjaanByAdmin(s.db, alumniID, adminID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error(), "success": false})
	}

	return c.JSON(fiber.Map{
		"message":       fmt.Sprintf("Berhasil menghapus %d pekerjaan milik alumni ID %d", count, alumniID),
		"success":       true,
		"deleted_count": count,
	})
}

// 🧱 Ambil data yang disoftdelete
func (s *JobService) GetTrashed(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := int(c.Locals("user_id").(float64))

	var jobs []model.Pekerjaan
	var err error

	if role == "admin" {
		jobs, err = repository.GetTrashedJobs(s.db)
	} else {
		jobs, err = repository.GetTrashedJobsByAlumni(s.db, userID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error(), "success": false})
	}
	if len(jobs) == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "Tidak ada data yang dihapus", "success": false})
	}

	return c.JSON(fiber.Map{"message": "Data yang dihapus berhasil diambil", "data": jobs, "success": true})
}

// 🧱 Restore pekerjaan
func (s *JobService) Restore(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	role := c.Locals("role").(string)
	userID := int(c.Locals("user_id").(float64))

	var err error
	if role == "admin" {
		err = repository.RestoreJobs(s.db, id)
	} else {
		err = repository.RestoreJobsByAlumni(s.db, id, userID)
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error(), "success": false})
	}
	return c.JSON(fiber.Map{"message": "Data berhasil direstore", "success": true})
}

// 🧱 Hard delete pekerjaan
func (s *JobService) HardDelete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	role := c.Locals("role").(string)
	userID := int(c.Locals("user_id").(float64))

	var err error
	if role == "admin" {
		err = repository.HardDeleteJob(s.db, id)
	} else {
		err = repository.HardDeleteJobByAlumni(s.db, id, userID)
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error(), "success": false})
	}
	return c.JSON(fiber.Map{"message": "Data berhasil dihapus permanen", "success": true})
}