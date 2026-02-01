package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database Conneced Successfully")
	return db, nil
}

// package database

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func InitDB(connectionString string) (*pgxpool.Pool, error) {
// 	// Konfigurasi pool agar lebih stabil di Supabase
// 	config, err := pgxpool.ParseConfig(connectionString)
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal parse config: %v", err)
// 	}

// 	// Setting tambahan agar tidak mudah EOF
// 	config.MaxConns = 10
// 	config.MinConns = 2
// 	config.MaxConnIdleTime = 5 * time.Minute

// 	// Membuka koneksi
// 	db, err := pgxpool.NewWithConfig(context.Background(), config)
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal membuat pool: %v", err)
// 	}

// 	// Tes koneksi
// 	err = db.Ping(context.Background())
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal ping database: %v", err)
// 	}

// 	log.Println("Database Connected Successfully (using PGX)")
// 	return db, nil
// }
