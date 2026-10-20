package router

import (
	"taskflow-api/internal/auth"
	"taskflow-api/internal/handler"
	"taskflow-api/internal/middleware"
	"taskflow-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RouterConfig struct {
	App          *fiber.App
	AuthHandler  *handler.AuthHandler
	TokenManager auth.TokenManager
	RedisClient  *redis.Client
}

func SetupRouter(config *RouterConfig) {
	config.App.Get("/health", func(c *fiber.Ctx) error {
		return utils.SuccessResponse(c, fiber.StatusOK, "TF-00-001", "Server is healthy", fiber.Map{
			"status": "up",
		})
	})

	api := config.App.Group("/api/v1")
	authGroup := api.Group("/auth")

	// Endpoint Publik
	authGroup.Post("/register", config.AuthHandler.Register)
	authGroup.Post("login", config.AuthHandler.Login)

	// Inisialisasi Middleware Auth
	authMiddleware := middleware.DeserializeUserMiddlware(config.TokenManager)

	// Endpoint Terproteksi (Hanya bisa diakses jika bawa Bearer Token valid)
	authGroup.Get("/me", authMiddleware, func(c *fiber.Ctx) error {
		// Mengambil data yang tadi diselipkan oleh middleware lewat c.Locals
		userID := c.Locals("user_id")
		userEmail := c.Locals("user_email")
		userRole := c.Locals("user_role")

		return utils.SuccessResponse(c, fiber.StatusOK, "TF-AUTH-007", "Profil user berhasil diambil", fiber.Map{
			"id":   userID,
			"email": userEmail,
			"role": userRole,
		})
	})
}
