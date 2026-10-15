package utils

import (
	"github.com/gofiber/fiber/v3"
)

// Format JSON untuk respons sukses.
type ResponseSuccess struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Format JSON untuk respons gagal.
type ResponseError struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse mengirim respons sukses dengan status code, pesan, dan data.
func SuccessResponse(c fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(ResponseSuccess{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirim respons gagal dengan status code, pesan, dan detail error.
func ErrorResponse(c fiber.Ctx, statusCode int, message string, errors interface{}) error {
	return c.Status(statusCode).JSON(ResponseError{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
