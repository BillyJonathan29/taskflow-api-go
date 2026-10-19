package utils

import "github.com/gofiber/fiber/v2"

type ResponseFormat[T any] struct {
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

// SuccessResponse mengembalikan format standar dengan generic data
func SuccessResponse[T any](c *fiber.Ctx, httpStatusCode int, customCode string, message string, data T) error {
	return c.Status(httpStatusCode).JSON(ResponseFormat[T]{
		StatusCode: customCode,
		Message:    message,
		Data:       data,
	})
}

// ErrorResponse untuk respon gagal tanpa data
func ErrorResponse(c *fiber.Ctx, httpStatusCode int, customCode string, message string) error {
	return c.Status(httpStatusCode).JSON(ResponseFormat[any]{
		StatusCode: customCode,
		Message:    message,
		Data:       nil,
	})
}

// // Format JSON untuk respons sukses.
// type ResponseSuccess struct {

// 	Success bool        `json:"success"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data,omitempty"`
// }

// // Format JSON untuk respons gagal.
// type ResponseError struct {
// 	Success bool        `json:"success"`
// 	Message string      `json:"message"`
// 	Errors  interface{} `json:"errors,omitempty"`
// }

// // SuccessResponse mengirim respons sukses dengan status code, pesan, dan data.
// func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
// 	return c.Status(statusCode).JSON(ResponseSuccess{
// 		Success: true,
// 		Message: message,
// 		Data:    data,
// 	})
// }

// // ErrorResponse mengirim respons gagal dengan status code, pesan, dan detail error.
// func ErrorResponse(c *fiber.Ctx, statusCode int, message string, errors interface{}) error {
// 	return c.Status(statusCode).JSON(ResponseError{
// 		Success: false,
// 		Message: message,
// 		Errors:  errors,
// 	})
// }
