package main

import (
	"fmt"
	"log"
	"taskflow-api/internal/config"
	"taskflow-api/internal/domain"
	"taskflow-api/pkg/database"
	"taskflow-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func main() {
	// 1. Inisialisasi Konfigurasi
	cfg := config.LoadConfig()

	// 2. inisialisasi gorm
	db := database.NewPostgresGORM(cfg)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan sqlDB: %v", err)
	}
	defer sqlDB.Close()

	// 3. Auto Migrate (GORM akan otomatis membuat tabel 'users' jika belum ada)
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("Gagal melakukan migrasi table user: %v", err)
	}
	log.Println("Migrasi ke database berhasil")

	// 4.inisialisasi fiber
	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		dbStatus := "connected"
		if err := sqlDB.Ping(); err != nil {
			dbStatus = "disconected"
		}

		return utils.SuccessResponse(c, fiber.StatusOK, "Server is healthy", fiber.Map{
			"environment": cfg.AppEnv,
			"port":        cfg.AppPort,
			"database":    dbStatus,
		})
	})

	listerAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server berjalan di port %s (%s)", cfg.AppPort, cfg.AppEnv)
	log.Fatal(app.Listen(listerAddr))
}
