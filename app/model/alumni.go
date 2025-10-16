package model

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type Alumni struct {
	ID         int       `json:"id"`
	NIM        string    `json:"nim"`
	Nama       string    `json:"nama"`
	Jurusan    string    `json:"jurusan"`
	Angkatan   int       `json:"angkatan"`
	TahunLulus int       `json:"tahun_lulus"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	NoTelp     *string   `json:"no_telepon"`
	Alamat     *string   `json:"alamat"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Alumni Alumni `json:"alumni"`
	Token  string `json:"token"`
}

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims

}
