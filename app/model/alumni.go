package model

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Alumni untuk PostgreSQL & MongoDB
type Alumni struct {
    MongoID     primitive.ObjectID `json:"id" bson:"_id,omitempty" db:"mongo_id"`
    ID          int                `json:"id" db:"id"`
    NIM         string             `json:"nim" bson:"nim" db:"nim"`
    Nama        string             `json:"nama" bson:"nama" db:"nama"`
    Jurusan     string             `json:"jurusan" bson:"jurusan" db:"jurusan"`
    Angkatan    string             `json:"angkatan" bson:"angkatan" db:"angkatan"`
    TahunLulus string              `json:"tahun_lulus" bson:"tahun_lulus" db:"tahun_lulus"`
    Email       string             `json:"email" bson:"email" db:"email"`
    Password    string             `json:"password" bson:"password" db:"password"`
    NoTelp      *string            `json:"no_telp,omitempty" bson:"no_telp,omitempty" db:"no_telp"`
    Alamat      *string            `json:"alamat,omitempty" bson:"alamat,omitempty" db:"alamat"`
    Role        string             `json:"role" bson:"role" db:"role"`
    CreatedAt   time.Time          `json:"created_at" bson:"created_at" db:"created_at"`
    UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

// Struct untuk login request
type LoginRequest struct {
    Email    string `json:"email" bson:"email"`
    Password string `json:"password" bson:"password"`
}

// Struct untuk response login
type LoginResponse struct {
    Alumni Alumni `json:"alumni"`
    Token  string `json:"token"`
}

// JWT Claims
type JWTClaims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
