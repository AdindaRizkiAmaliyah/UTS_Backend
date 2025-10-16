package model

import "time"

type Pekerjaan struct {
	ID                  int        `json:"id" db:"id"` // serial
	AlumniID            int        `json:"alumni_id" db:"alumni_id"` // integer (foreign key)
	NamaPerusahaan      string     `json:"nama_perusahaan" db:"nama_perusahaan"` // varchar(100)
	PosisiJabatan       string     `json:"posisi_jabatan" db:"posisi_jabatan"` // varchar(100)
	BidangIndustri      string     `json:"bidang_industri" db:"bidang_industri"` // varchar(100)
	LokasiKerja         string     `json:"lokasi_kerja" db:"lokasi_kerja"` // varchar(100)
	GajiRange           *string    `json:"gaji_range,omitempty" db:"gaji_range"` // varchar(50)
	TanggalMulaiKerja   *time.Time `json:"tanggal_mulai_kerja,omitempty" db:"tanggal_mulai_kerja"` // date
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja,omitempty" db:"tanggal_selesai_kerja"` // date
	StatusPekerjaan     *string    `json:"status_pekerjaan,omitempty" db:"status_pekerjaan"` // varchar(50)
	DeskripsiPekerjaan  *string    `json:"deskripsi_pekerjaan,omitempty" db:"deskripsi_pekerjaan"` // text
	CreatedAt           time.Time  `json:"created_at" db:"created_at"` // timestamp
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"` // timestamp
	IsDeleted           *time.Time `json:"is_deleted,omitempty" db:"is_deleted"` // timestamp
	DeletedBy           *string    `json:"deleted_by,omitempty" db:"deleted_by"` // varchar(50)
}
