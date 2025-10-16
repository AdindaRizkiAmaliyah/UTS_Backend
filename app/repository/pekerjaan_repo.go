package repository

import (
	"clean-archi/app/model"
	"database/sql"
	"fmt"
	"strings"
)

// GetAllJobs returns all non-deleted pekerjaan (is_deleted IS NULL)
func GetAllJobs(db *sql.DB) ([]model.Pekerjaan, error) {
	rows, err := db.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
		       created_at, updated_at, deleted_by, is_deleted
		FROM pekerjaan
		WHERE is_deleted IS NULL
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
			&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedBy, &p.IsDeleted,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// GetJobByID returns single pekerjaan by id
func GetJobByID(db *sql.DB, id int) (*model.Pekerjaan, error) {
	row := db.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
		       created_at, updated_at, deleted_by, is_deleted
		FROM pekerjaan
		WHERE id = $1
	`, id)

	var p model.Pekerjaan
	err := row.Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
		&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
		&p.CreatedAt, &p.UpdatedAt, &p.DeletedBy, &p.IsDeleted,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// GetJobsByAlumniID returns all pekerjaan by alumni (not deleted)
func GetJobsByAlumniID(db *sql.DB, alumniID int) ([]model.Pekerjaan, error) {
	rows, err := db.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
		       created_at, updated_at, deleted_by, is_deleted
		FROM pekerjaan
		WHERE alumni_id = $1 AND is_deleted IS NULL
		ORDER BY created_at DESC
	`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
			&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedBy, &p.IsDeleted,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// GetAllJobsWithFilter - filter dinamis
func GetAllJobsWithFilter(db *sql.DB, filters map[string]string) ([]model.Pekerjaan, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
		       created_at, updated_at, deleted_by, is_deleted
		FROM pekerjaan
		WHERE is_deleted IS NULL
	`

	args := []interface{}{}
	conditions := []string{}

	if v, ok := filters["nama_perusahaan"]; ok && v != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(nama_perusahaan) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+v+"%")
	}
	if v, ok := filters["posisi_jabatan"]; ok && v != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(posisi_jabatan) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+v+"%")
	}
	if v, ok := filters["bidang_industri"]; ok && v != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(bidang_industri) LIKE LOWER($%d)", len(args)+1))
		args = append(args, "%"+v+"%")
	}
	if v, ok := filters["status_pekerjaan"]; ok && v != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(status_pekerjaan) = LOWER($%d)", len(args)+1))
		args = append(args, v)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
			&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedBy, &p.IsDeleted,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// CreateJob - insert baru
func CreateJob(db *sql.DB, pekerjaan *model.Pekerjaan) (*model.Pekerjaan, error) {
	query := `
		INSERT INTO pekerjaan
			(alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
			 tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
			 created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
		RETURNING id, created_at, updated_at
	`
	err := db.QueryRow(query,
		pekerjaan.AlumniID, pekerjaan.NamaPerusahaan, pekerjaan.PosisiJabatan,
		pekerjaan.BidangIndustri, pekerjaan.LokasiKerja, pekerjaan.GajiRange,
		pekerjaan.TanggalMulaiKerja, pekerjaan.TanggalSelesaiKerja,
		pekerjaan.StatusPekerjaan, pekerjaan.DeskripsiPekerjaan,
	).Scan(&pekerjaan.ID, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return pekerjaan, nil
}

// UpdateJob - update data pekerjaan
func UpdateJob(db *sql.DB, id int, pekerjaan *model.Pekerjaan) (*model.Pekerjaan, error) {
	query := `
		UPDATE pekerjaan
		SET 
			nama_perusahaan = $1,
			posisi_jabatan = $2,
			bidang_industri = $3,
			lokasi_kerja = $4,
			gaji_range = $5,
			tanggal_mulai_kerja = $6,
			tanggal_selesai_kerja = $7,
			status_pekerjaan = $8,
			deskripsi_pekerjaan = $9,
			updated_at = NOW()
		WHERE id = $10 AND is_deleted IS NULL
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
		          gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
		          deskripsi_pekerjaan, created_at, updated_at, deleted_by, is_deleted
	`

	var updated model.Pekerjaan
	err := db.QueryRow(
		query,
		pekerjaan.NamaPerusahaan,
		pekerjaan.PosisiJabatan,
		pekerjaan.BidangIndustri,
		pekerjaan.LokasiKerja,
		pekerjaan.GajiRange,
		pekerjaan.TanggalMulaiKerja,
		pekerjaan.TanggalSelesaiKerja,
		pekerjaan.StatusPekerjaan,
		pekerjaan.DeskripsiPekerjaan,
		id,
	).Scan(
		&updated.ID, &updated.AlumniID, &updated.NamaPerusahaan, &updated.PosisiJabatan,
		&updated.BidangIndustri, &updated.LokasiKerja, &updated.GajiRange,
		&updated.TanggalMulaiKerja, &updated.TanggalSelesaiKerja, &updated.StatusPekerjaan,
		&updated.DeskripsiPekerjaan, &updated.CreatedAt, &updated.UpdatedAt,
		&updated.DeletedBy, &updated.IsDeleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("data tidak ditemukan atau sudah dihapus")
		}
		return nil, err
	}

	return &updated, nil
}


// Soft delete satu pekerjaan oleh alumni
func SoftDeletePekerjaanByAlumni(db *sql.DB, jobID, alumniID int) error {
	query := `
		UPDATE pekerjaan
		SET is_deleted = NOW(),
		    deleted_by = $2,
		    updated_at = NOW()
		WHERE id = $1 AND alumni_id = $2 AND is_deleted IS NULL
	`
	res, err := db.Exec(query, jobID, alumniID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("anda tidak memiliki izin untuk menghapus pekerjaan ini atau data tidak ditemukan")
	}

	return nil
}

// Soft delete semua pekerjaan milik alumni (oleh admin)
func SoftDeleteAllPekerjaanByAdmin(db *sql.DB, alumniID, adminID int) (int64, error) {
	query := `
		UPDATE pekerjaan
		SET is_deleted = NOW(),
		    deleted_by = $2,
		    updated_at = NOW()
		WHERE alumni_id = $1 AND is_deleted IS NULL
	`
	res, err := db.Exec(query, alumniID, adminID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Ambil semua data yang disoftdelete (khusus admin)
func GetTrashedJobs(db *sql.DB) ([]model.Pekerjaan, error) {
	rows, err := db.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
		       created_at, updated_at, deleted_by, is_deleted
		FROM pekerjaan
		WHERE is_deleted IS NOT NULL
		ORDER BY is_deleted DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []model.Pekerjaan
	for rows.Next() {
		var j model.Pekerjaan
		err := rows.Scan(
			&j.ID, &j.AlumniID, &j.NamaPerusahaan, &j.PosisiJabatan, &j.BidangIndustri,
			&j.LokasiKerja, &j.GajiRange, &j.TanggalMulaiKerja, &j.TanggalSelesaiKerja,
			&j.StatusPekerjaan, &j.DeskripsiPekerjaan, &j.CreatedAt, &j.UpdatedAt,
			&j.DeletedBy, &j.IsDeleted,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// Ambil pekerjaan yang dihapus milik alumni tertentu
func GetTrashedJobsByAlumni(db *sql.DB, alumniID int) ([]model.Pekerjaan, error) {
	rows, err := db.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		       tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
		       created_at, updated_at, deleted_by, is_deleted
		FROM pekerjaan
		WHERE is_deleted IS NOT NULL AND alumni_id = $1
		ORDER BY is_deleted DESC
	`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []model.Pekerjaan
	for rows.Next() {
		var j model.Pekerjaan
		err := rows.Scan(
			&j.ID, &j.AlumniID, &j.NamaPerusahaan, &j.PosisiJabatan, &j.BidangIndustri,
			&j.LokasiKerja, &j.GajiRange, &j.TanggalMulaiKerja, &j.TanggalSelesaiKerja,
			&j.StatusPekerjaan, &j.DeskripsiPekerjaan, &j.CreatedAt, &j.UpdatedAt,
			&j.DeletedBy, &j.IsDeleted,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// Restore pekerjaan (admin)
func RestoreJobs(db *sql.DB, jobID int) error {
	_, err := db.Exec(`
		UPDATE pekerjaan
		SET is_deleted = NULL,
		    deleted_by = NULL,
		    updated_at = NOW()
		WHERE id = $1
	`, jobID)
	return err
}

// Restore pekerjaan milik alumni sendiri
func RestoreJobsByAlumni(db *sql.DB, jobID, alumniID int) error {
	res, err := db.Exec(`
		UPDATE pekerjaan
		SET is_deleted = NULL, deleted_by = NULL, updated_at = NOW()
		WHERE id = $1 AND alumni_id = $2 AND is_deleted IS NOT NULL
	`, jobID, alumniID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("data tidak ditemukan atau bukan milik anda")
	}
	return nil
}

// Hard delete (admin)
func HardDeleteJob(db *sql.DB, jobID int) error {
	_, err := db.Exec(`DELETE FROM pekerjaan WHERE id = $1`, jobID)
	return err
}

// Hard delete milik alumni sendiri
func HardDeleteJobByAlumni(db *sql.DB, jobID, alumniID int) error {
	res, err := db.Exec(`DELETE FROM pekerjaan WHERE id = $1 AND alumni_id = $2`, jobID, alumniID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("data tidak ditemukan atau bukan milik anda")
	}
	return nil
}
