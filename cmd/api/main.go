package main

import (
	"log"
	"os"
	"taskflow-api/internal/bootstrap"
	"taskflow-api/internal/config"
	"taskflow-api/internal/model"
	"taskflow-api/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load Environment Variables
	config.LoadEnv()

	// 2. Init Database & Redis
	db := config.InitDB()
	redisClient := config.InitRedis()

	// Auto Migrate DB
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("AutoMigrate Failed: %v", err)
	}
	// 3. Wiring Dependencies lewat Bootstrap
	handlers := bootstrap.InitApp(db, redisClient)
	// 4. Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		AppName: "TaskFlow Appi Core",
	})

	// 5. Setup Router
	router.SetupRouter(&router.RouterConfig{
		App:          app,
		AuthHandler:  handlers.AuthHandler,
		TokenManager: handlers.TokenManager,
		RedisClient:  redisClient,
	})

	// 6. Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Server Aktif di http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
