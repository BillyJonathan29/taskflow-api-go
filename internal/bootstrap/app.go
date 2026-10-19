package bootstrap

import (
	"taskflow-api/internal/handler"
	"taskflow-api/internal/repository"
	"taskflow-api/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	AuthHandler *handler.AuthHandler
}

func InitApp(db *gorm.DB, redisClient *redis.Client) *Handler {
	// 1. inisialisasi repo
	userRepo := repository.NewUserRepository(db)

	// 2. inisialisasi Repo
	authService := service.NewAuthService(userRepo)

	// 3. inisialisas Handler
	authHandler := handler.NewAuthHandler(authService)

	return &Handler{
		AuthHandler: authHandler,
	}
}
