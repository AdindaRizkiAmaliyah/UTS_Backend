// package service

// import (
// 	"database/sql"
// 	"os"
// 	"github.com/gofiber/fiber/v2"
// 	"clean-archi/app/repository"
// 	"clean-archi/app/model"
// 	"strconv"
// 	"strings"
// )

// type AlumniService struct {
// 	db *sql.DB
// }

// func NewAlumniService(db *sql.DB) *AlumniService {
// 	return &AlumniService{db: db}
// }

// // validateAPIKey adalah helper function untuk validasi API key
// func (s *AlumniService) validateAPIKey(c *fiber.Ctx) error {
// 	// Cek dari header X-API-KEY terlebih dahulu
// 	key := c.Get("X-API-KEY")
// 	if key == "" {
// 		// Jika tidak ada di header, cek dari URL params
// 		key = c.Params("key")
// 	}
	
// 	if key != os.Getenv("API_KEY") {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"message": "Key tidak valid",
// 			"success": false,
// 		})
// 	}
// 	return nil
// }

// // CheckAlumni - Mengecek apakah mahasiswa adalah alumni
// func (s *AlumniService) CheckAlumni(c *fiber.Ctx) error {
// 	if err := s.validateAPIKey(c); err != nil {
// 		return err
// 	}

// 	nim := c.FormValue("nim")
// 	if nim == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "NIM wajib diisi",
// 			"success": false,
// 		})
// 	}

// 	alumni, err := repository.CheckAlumniByNim(s.db, nim)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 				"message":  "Mahasiswa bukan alumni",
// 				"success":  true,
// 				"isAlumni": false,
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Gagal cek alumni karena " + err.Error(),
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message":  "Berhasil mendapatkan data alumni",
// 		"success":  true,
// 		"isAlumni": true,
// 		"alumni":   alumni,
// 	})
// }

// // GetAll - Mendapatkan semua data alumni dengan pagination
// func (s *AlumniService) GetAll(c *fiber.Ctx) error {
// 	if err := s.validateAPIKey(c); err != nil {
// 		return err
// 	}

// 	// Ambil query params
// 	page, _ := strconv.Atoi(c.Query("page", "1"))
// 	limit, _ := strconv.Atoi(c.Query("limit", "10"))
// 	sortBy := c.Query("sortBy", "id")
// 	order := c.Query("order", "asc")
// 	search := c.Query("search", "")

// 	offset := (page - 1) * limit

// 	// Whitelist kolom agar aman dari SQL Injection
// 	whitelist := map[string]bool{
// 		"id": true, "nim": true, "nama": true, "jurusan": true, 
// 		"angkatan": true, "tahun_lulus": true, "email": true, "created_at": true,
// 	}
// 	if !whitelist[sortBy] {
// 		sortBy = "id"
// 	}
// 	if strings.ToLower(order) != "desc" {
// 		order = "asc"
// 	}

// 	alumniList, err := repository.GetAlumniPaginated(s.db, search, sortBy, order, limit, offset)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error: " + err.Error(),
// 			"success": false,
// 		})
// 	}

// 	total, err := repository.CountAlumni(s.db, search)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error menghitung total: " + err.Error(),
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Berhasil mendapatkan data alumni",
// 		"success": true,
// 		"data":    alumniList,
// 		"meta": fiber.Map{
// 			"page":   page,
// 			"limit":  limit,
// 			"total":  total,
// 			"pages":  (total + limit - 1) / limit,
// 			"sortBy": sortBy,
// 			"order":  order,
// 			"search": search,
// 		},
// 	})
// }

// // GetByID - Mendapatkan data alumni berdasarkan ID
// func (s *AlumniService) GetByID(c *fiber.Ctx) error {
// 	if err := s.validateAPIKey(c); err != nil {
// 		return err
// 	}

// 	id, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "ID tidak valid",
// 			"success": false,
// 		})
// 	}

// 	alumni, err := repository.GetAlumniByID(s.db, id)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Error: " + err.Error(),
// 			"success": false,
// 		})
// 	}

// 	if alumni == nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"message": "Alumni tidak ditemukan",
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Berhasil mendapatkan data alumni",
// 		"success": true,
// 		"alumni": alumni,
// 	})
// }

// // Create - Membuat data alumni baru
// func (s *AlumniService) Create(c *fiber.Ctx) error {
// 	if err := s.validateAPIKey(c); err != nil {
// 		return err
// 	}

// 	var alumni model.Alumni
// 	if err := c.BodyParser(&alumni); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Request body tidak valid",
// 			"success": false,
// 		})
// 	}

// 	// Insert ke DB
// 	savedAlumni, err := repository.CreateAlumni(s.db, &alumni)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Gagal menambahkan alumni: " + err.Error(),
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "Alumni berhasil ditambahkan",
// 		"success": true,
// 		"alumni":  savedAlumni,
// 	})
// }

// // Update - Mengupdate data alumni
// func (s *AlumniService) Update(c *fiber.Ctx) error {
// 	if err := s.validateAPIKey(c); err != nil {
// 		return err
// 	}

// 	// Ambil ID dari URL
// 	id := c.Params("id")

// 	// Parse body ke struct Alumni
// 	var alumni model.Alumni
// 	if err := c.BodyParser(&alumni); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Request body tidak valid",
// 			"success": false,
// 		})
// 	}

// 	// Update ke DB
// 	updatedAlumni, err := repository.UpdateAlumni(s.db, id, &alumni)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Gagal mengupdate alumni: " + err.Error(),
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Alumni berhasil diperbarui",
// 		"success": true,
// 		"alumni": updatedAlumni,
// 	})
// }

// // Delete - Menghapus data alumni
// func (s *AlumniService) Delete(c *fiber.Ctx) error {
// 	if err := s.validateAPIKey(c); err != nil {
// 		return err
// 	}

// 	// Ambil ID dari URL
// 	id := c.Params("id")

// 	// Hapus alumni dari DB
// 	err := repository.DeleteAlumni(s.db, id)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Gagal menghapus alumni: " + err.Error(),
// 			"success": false,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"message": "Alumni berhasil dihapus",
// 		"success": true,
// 	})
// }

// // Backward compatibility functions (opsional, jika masih diperlukan)
// func CheckAlumniService(c *fiber.Ctx, db *sql.DB) error {
// 	service := NewAlumniService(db)
// 	return service.CheckAlumni(c)
// }

// func GetAllAlumniService(c *fiber.Ctx, db *sql.DB) error {
// 	service := NewAlumniService(db)
// 	return service.GetAll(c)
// }

// func GetAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
// 	service := NewAlumniService(db)
// 	return service.GetByID(c)
// }

// func CreateAlumniService(c *fiber.Ctx, db *sql.DB) error {
// 	service := NewAlumniService(db)
// 	return service.Create(c)
// }

// func UpdateAlumniService(c *fiber.Ctx, db *sql.DB) error {
// 	service := NewAlumniService(db)
// 	return service.Update(c)
// }

// func DeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
// 	service := NewAlumniService(db)
// 	return service.Delete(c)
// }

package service

import (
    "clean-archi/app/model"
    "clean-archi/app/repository"
    "context"
    "os"
    "time"

    "github.com/gofiber/fiber/v2"
)

type AlumniService struct {
    repo repository.AlumniRepository
}

// Constructor — sekarang menerima interface, bukan MongoRepo
func NewAlumniService(repo repository.AlumniRepository) *AlumniService {
    return &AlumniService{repo: repo}
}

// ========================= BUSINESS LOGIC (UNTUK TESTING) =========================

func (s *AlumniService) CreateAlumni(ctx context.Context, a *model.Alumni) (*model.Alumni, error) {
    return s.repo.Create(ctx, a)
}

func (s *AlumniService) GetAlumniByID(ctx context.Context, id string) (*model.Alumni, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *AlumniService) UpdateAlumni(ctx context.Context, id string, a *model.Alumni) error {
    return s.repo.Update(ctx, id, a)
}

func (s *AlumniService) DeleteAlumni(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}

// ========================= VALIDASI API KEY =========================

func (s *AlumniService) validateAPIKey(c *fiber.Ctx) error {
    key := c.Get("X-API-KEY")
    if key == "" {
        key = c.Query("key")
    }

    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "API Key tidak valid",
            "success": false,
        })
    }
    return nil
}

// ========================= HTTP HANDLERS (FIBER) =========================

// GET ALL
func (s *AlumniService) GetAll(c *fiber.Ctx) error {
    if err := s.validateAPIKey(c); err != nil {
        return err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    alumniList, err := s.repo.GetAll(ctx)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mengambil data: " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil mendapatkan semua alumni",
        "success": true,
        "data":    alumniList,
    })
}

// GET BY ID
func (s *AlumniService) GetByID(c *fiber.Ctx) error {
    if err := s.validateAPIKey(c); err != nil {
        return err
    }

    id := c.Params("id")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    alumni, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "Alumni tidak ditemukan",
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Berhasil mendapatkan data alumni",
        "success": true,
        "data":    alumni,
    })
}

// CREATE
func (s *AlumniService) Create(c *fiber.Ctx) error {
    if err := s.validateAPIKey(c); err != nil {
        return err
    }

    var alumni model.Alumni
    if err := c.BodyParser(&alumni); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Body request tidak valid",
            "success": false,
        })
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    newAlumni, err := s.repo.Create(ctx, &alumni)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menambahkan alumni: " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Alumni berhasil ditambahkan",
        "success": true,
        "data":    newAlumni,
    })
}

// UPDATE
func (s *AlumniService) Update(c *fiber.Ctx) error {
    if err := s.validateAPIKey(c); err != nil {
        return err
    }

    id := c.Params("id")
    var alumni model.Alumni

    if err := c.BodyParser(&alumni); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Body request tidak valid",
            "success": false,
        })
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := s.repo.Update(ctx, id, &alumni); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal memperbarui alumni: " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Alumni berhasil diperbarui",
        "success": true,
    })
}

// DELETE
func (s *AlumniService) Delete(c *fiber.Ctx) error {
    if err := s.validateAPIKey(c); err != nil {
        return err
    }

    id := c.Params("id")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := s.repo.Delete(ctx, id); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal menghapus alumni: " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Alumni berhasil dihapus",
        "success": true,
    })
}
