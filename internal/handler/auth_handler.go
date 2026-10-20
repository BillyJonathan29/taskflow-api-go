package handler

import (
	"taskflow-api/internal/dto"
	"taskflow-api/internal/service"
	"taskflow-api/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// 1. Parse JSON body
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "TF-01-003", "Format JSON tidak valid")
	}
	// 2. Validasi input
	if err := h.validate.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "TF-01-002", err.Error())
	}

	res, err := h.authService.Register(req)
	if err != nil {
		if err.Error() == "email sudah terdaftar" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "TF-01-004", err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "TF-01-500", err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, "TF-01-001", "User berhasil didaftarkan", res)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// 1. Parsing body JSON
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "TF-01-003", "Format JSON tidak valid")
	}

	// 2. Validasi format email & required
	if err := h.validate.Struct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "TF-01-002", err.Error())
	}

	// 3. Panggil service login dengan membawa context Fiber (c.Context())
	res, err := h.authService.Login(c.Context(), req)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "TF-01-005", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "TF-01-006", "Login berhasil", res)

}
