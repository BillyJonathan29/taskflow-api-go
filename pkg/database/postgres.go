package database

import (
	"fmt"
	"log"
	"time"
	"taskflow-api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresGORM(cfg *config.Config) *gorm.DB {
	// 1. Format DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	// 2. Tentukan level log GORM (Info di dev agar terlihat SQL query-nya)
	gormLogLevel := logger.Warn
	if cfg.AppEnv == "development" {
		gormLogLevel = logger.Info
	}

	// inisialisasi gorm
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil{
		log.Fatalf("Gagal terhubung ke postgres lewat GORM: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil{
		log.Fatalf("Gagal mengambil instance sql.DB dari GORM: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)                 // Maksimal 25 koneksi aktif
	sqlDB.SetMaxIdleConns(10)                 // Minimal 10 koneksi siaga
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Usia maksimal koneksi
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // Batas waktu koneksi idle

	if err := sqlDB.Ping(); err != nil{
		log.Fatalf("Ping ke database gagal: %v", err)
	}

	log.Println("Berhasil terhubung ke PostgreSQL via gorm")
	return db
}
