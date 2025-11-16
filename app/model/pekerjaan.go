package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =======================================
// Model Pekerjaan (Bisa dimiliki >1 alumni)
// =======================================
type Pekerjaan struct {
	MongoID            primitive.ObjectID   `bson:"_id,omitempty" json:"mongo_id,omitempty"`
	KpnID              int                  `json:"id_kpn,omitempty" bson:"id_kpn,omitempty"`
	AlumniIDs          []primitive.ObjectID `json:"alumni_id" bson:"alumni_id"` // <-- perbaikan: array of ObjectId
	NamaPerusahaan     string               `json:"nama_perusahaan" bson:"nama_perusahaan"`
	PosisiJabatan      string               `json:"posisi_jabatan" bson:"posisi_jabatan"`
	BidangIndustri     string               `json:"bidang_industri" bson:"bidang_industri"`
	LokasiKerja        string               `json:"lokasi_kerja" bson:"lokasi_kerja"`
	DeskripsiPekerjaan string               `json:"deskripsi_pekerjaan" bson:"deskripsi_pekerjaan"`
	TanggalMulaiKerja  time.Time            `json:"tanggal_mulai_kerja" bson:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja time.Time           `json:"tanggal_selesai_kerja,omitempty" bson:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan    string               `json:"status_pekerjaan,omitempty" bson:"status_pekerjaan,omitempty"`
	CreatedAt          time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at" bson:"updated_at"`
	IsDeleted          bool                 `json:"is_deleted" bson:"is_deleted"`
	DeletedBy          *string              `json:"deleted_by,omitempty" bson:"deleted_by,omitempty"`
}

// =======================================
// DTO untuk input lewat Swagger
// =======================================
type CreatePekerjaanRequest struct {
	AlumniIDs          []string  `json:"alumni_id"` // array of string ObjectID dari Swagger
	NamaPerusahaan     string    `json:"nama_perusahaan"`
	PosisiJabatan      string    `json:"posisi_jabatan"`
	BidangIndustri     string    `json:"bidang_industri"`
	LokasiKerja        string    `json:"lokasi_kerja"`
	DeskripsiPekerjaan string    `json:"deskripsi_pekerjaan"`
	TanggalMulaiKerja  time.Time `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan    string    `json:"status_pekerjaan"`
}
