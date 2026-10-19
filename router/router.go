package router

import (
	"taskflow-api/internal/handler"
	"taskflow-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RouterConfig struct {
	App         *fiber.App
	AuthHandler *handler.AuthHandler
	RedisClient *redis.Client
}

func SetupRouter(config *RouterConfig) {
	config.App.Get("/health", func(c *fiber.Ctx) error {
		return utils.SuccessResponse(c, fiber.StatusOK, "TF-00-001", "Server is healthy", fiber.Map{
			"status": "up",
		})
	})

	api := config.App.Group("/api/v1")
	auth := api.Group("/auth")
	auth.Post("/register", config.AuthHandler.Register)
}
