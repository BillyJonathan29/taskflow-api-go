package main

import (
	"fmt"
	"log"
	"taskflow-api/internal/config"
	delivery "taskflow-api/internal/delivery/http"
	"taskflow-api/internal/domain"
	"taskflow-api/internal/repository"
	"taskflow-api/internal/service"
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
	// log.Println("Migrasi ke database berhasil")

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := delivery.NewAuthHandler(authService)

	// 4.inisialisasi fiber
	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {

		return utils.SuccessResponse(c, fiber.StatusOK, "Server is healthy", fiber.Map{
			"environtment": cfg.AppEnv,
			"port":         cfg.AppPort,
		})
		// dbStatus := "connected"
		// if err := sqlDB.Ping(); err != nil {
		// 	dbStatus = "disconected"
		// }

		// return utils.SuccessResponse(c, fiber.StatusOK, "Server is healthy", fiber.Map{
		// 	"environment": cfg.AppEnv,
		// 	"port":        cfg.AppPort,
		// 	"database":    dbStatus,
		// })
	})
	v1 := app.Group("/api/v1")
	authRoutes := v1.Group("/auth")
	authRoutes.Post("/register", authHandler.Register)

	listerAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server berjalan di port %s (%s)", cfg.AppPort, cfg.AppEnv)
	log.Fatal(app.Listen(listerAddr))
}
