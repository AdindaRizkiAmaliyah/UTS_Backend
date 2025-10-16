package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	dsn := os.Getenv("DB_DSN")
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database tidak bisa diakses:", err)
	}

	log.Println("✅ Database terkoneksi")
}
