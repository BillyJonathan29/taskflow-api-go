package main

import (
	"fmt"
	"log"
	"taskflow-api/internal/config"
	"taskflow-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return utils.SuccessResponse(c, fiber.StatusOK, "Server is healthy", fiber.Map{
			"environment": cfg.AppEnv,
			"port":        cfg.AppPort,
		})
	})

	listerAddr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server berjalan di port %s (%s)", cfg.AppPort, cfg.AppEnv)
	log.Fatal(app.Listen(listerAddr))
}
