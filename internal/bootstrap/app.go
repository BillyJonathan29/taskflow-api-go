package bootstrap

import (
	"taskflow-api/internal/auth"
	"taskflow-api/internal/handler"
	"taskflow-api/internal/repository"
	"taskflow-api/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	AuthHandler  *handler.AuthHandler
	TokenManager auth.TokenManager // Export TokenManager untuk dipakai di middleware router
}

func InitApp(db *gorm.DB, redisClient *redis.Client) *Handler {
	// 1. inisialisasi repo
	userRepo := repository.NewUserRepository(db)

	// 2. Inisialisasi Token Manager (JWT + Redis)
	tokenManger := auth.NewTokenManager(redisClient)

	// 3. Inisialisasi Service (disuntikkan userRepo dan tokenManager)
	authService := service.NewAuthService(userRepo, tokenManger)

	// 4. Inisialisasi Handler
	authHandler := handler.NewAuthHandler(authService)
	return &Handler{
		AuthHandler: authHandler,
		TokenManager: tokenManger,
	}
}
