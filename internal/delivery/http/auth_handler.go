package http

import (
	"taskflow-api/internal/service"
	"taskflow-api/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req service.RegisterRequest

	// 1. Parse JSON body
	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format Json tidak valid", err.Error())
	}
	// 2. Validasi input
	if err := h.validate.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi gagal", err.Error())
	}
	// 3. Panggil business logic di service
	userResponse, err := h.authService.Register(req)
	if err != nil {
		if err.Error() == "Email sudah terdaftar" {
			return utils.ErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mendaftarkan user", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "User berhasil didaftarkan", userResponse)
}
