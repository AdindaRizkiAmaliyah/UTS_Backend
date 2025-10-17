package repository

import (
	"database/sql"
	"clean-archi/app/model"
    "fmt"
)

// CheckAlumniByNim mencari data alumni berdasarkan NIM
// Mengembalikan pointer ke model.Alumni jika ditemukan, atau error jika gagal
func CheckAlumniByNim(db *sql.DB, nim string) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	query := `SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
	FROM	alumni	WHERE	nim	=	$1	LIMIT	1`
	err := db.QueryRow(query, nim).Scan(&alumni.ID, &alumni.NIM, &alumni.Nama,
		&alumni.Jurusan, &alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.NoTelp, 
        &alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

// GetAllAlumni mengambil semua data alumni
// Mengembalikan slice alumniList dan error jika terjadi kesalahan
func GetAllAlumni(db *sql.DB) ([]model.Alumni, error) {
    rows, err := db.Query(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
    FROM alumni`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alumniList []model.Alumni
    for rows.Next() {
        var a model.Alumni
        err := rows.Scan( &a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email,  &a.NoTelp,  &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
        if err != nil {
            return nil, err
        }
        alumniList = append(alumniList, a)
    }

    return alumniList, nil
}

// GetAlumniByID mencari data alumni berdasarkan ID
// Mengembalikan pointer ke model.Alumni jika ditemukan, nil jika tidak ada, atau error jika gagal
func GetAlumniByID(db *sql.DB, id int) (*model.Alumni, error) {
    row := db.QueryRow(`SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, 
        email, no_telepon, alamat, created_at, updated_at 
        FROM alumni WHERE id=$1`, id)

    var alumni model.Alumni
    err := row.Scan(
        &alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
        &alumni.TahunLulus, &alumni.Email, &alumni.NoTelp, &alumni.Alamat,
        &alumni.CreatedAt, &alumni.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // data tidak ditemukan
        }
        return nil, err
    }

    return &alumni, nil
}

// CreateAlumni menambahkan data alumni baru ke database
// Mengembalikan pointer ke alumni yang baru dibuat, termasuk ID baru
func CreateAlumni(db *sql.DB, alumni *model.Alumni) (*model.Alumni, error) {
    query := `
        INSERT INTO alumni 
        (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
        RETURNING id
    `

    err := db.QueryRow(query,
        alumni.NIM,
        alumni.Nama,
        alumni.Jurusan,
        alumni.Angkatan,
        alumni.TahunLulus,
        alumni.Email,
        alumni.NoTelp,
        alumni.Alamat,
    ).Scan(&alumni.ID)

    if err != nil {
        return nil, err
    }

    return alumni, nil
}

// UpdateAlumni memperbarui data alumni berdasarkan ID
// Mengembalikan pointer ke data alumni yang sudah diperbarui
func UpdateAlumni(db *sql.DB, id string, alumni *model.Alumni) (*model.Alumni, error) {
    query := `
        UPDATE alumni
        SET nim = $1, nama = $2, jurusan = $3, angkatan = $4, 
            tahun_lulus = $5, email = $6, no_telepon = $7, alamat = $8, 
            updated_at = NOW()
        WHERE id = $9
        RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
    `

    row := db.QueryRow(query,
        alumni.NIM,
        alumni.Nama,
        alumni.Jurusan,
        alumni.Angkatan,
        alumni.TahunLulus,
        alumni.Email,
        alumni.NoTelp,
        alumni.Alamat,
        id,
    )

    var updated model.Alumni
    err := row.Scan(
        &updated.ID,
        &updated.NIM,
        &updated.Nama,
        &updated.Jurusan,
        &updated.Angkatan,
        &updated.TahunLulus,
        &updated.Email,
        &updated.NoTelp,
        &updated.Alamat,
        &updated.CreatedAt,
        &updated.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }

    return &updated, nil
}

// DeleteAlumni menghapus data alumni berdasarkan ID
// Mengembalikan error jika gagal atau jika data tidak ditemukan
func DeleteAlumni(db *sql.DB, id string) error {
    query := `DELETE FROM alumni WHERE id = $1`
    result, err := db.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("alumni dengan ID %s tidak ditemukan", id)
    }

    return nil
}

// GetAlumniPaginated mengambil data alumni dengan filter, urutkan, dan paginasi
// search = kata kunci pencarian, sortBy = kolom untuk diurutkan, order = ASC/DESC
// limit = jumlah data per halaman, offset = mulai dari data ke berapa
func GetAlumniPaginated(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
    query := fmt.Sprintf(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
        FROM alumni
        WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

    rows, err := db.Query(query, "%"+search+"%", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alumniList []model.Alumni
    for rows.Next() {
        var a model.Alumni
        err := rows.Scan(
            &a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
            &a.Email, &a.NoTelp, &a.Alamat, &a.CreatedAt, &a.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        alumniList = append(alumniList, a)
    }
    return alumniList, nil
}

// CountAlumni menghitung jumlah alumni sesuai kata kunci pencarian
// Berguna untuk pagination
func CountAlumni(db *sql.DB, search string) (int, error) {
    var total int
    query := `SELECT COUNT(*) FROM alumni 
              WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1`
    err := db.QueryRow(query, "%"+search+"%").Scan(&total)
    if err != nil {
        return 0, err
    }
    return total, nil
}











