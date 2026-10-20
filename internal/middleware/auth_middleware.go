package middleware

import (
	"strings"
	"taskflow-api/internal/auth"
	"taskflow-api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func DeserializeUserMiddlware(tokenManager auth.TokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil nilai header "Authorization"
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "TF-AUTH-001", "Header Authorization tidak ditemukan")
		}

		// 2. Format header harus berupa: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "TF-AUTH-002", "Format token tidak valid (wajib: Bearer <token>)")
		}

		tokenString := parts[1]

		// 3. Verifikasi signature JWT
		claims, err := tokenManager.VerifyToken(tokenString)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "TF-AUTH-003", "Token tidak valid atau sudah kadaluwarsa")
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "TF-AUTH-004", "Payload token tidak sah")
		}

		userID := uint(userIDFloat)

		// 4. Cek validitas sesi di Redis (Fitur Anti-Lag & Session Control)
		isValidSession := tokenManager.ValidateSession(c.Context(), userID, tokenString)
		if !isValidSession {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "TF-AUTH-005", "Sesi login telah berakhir di redis, silah login ulang")
		}

		// 5. Simpan data user ke Locals Fiber agar bisa dibaca oleh Handler berikutnya
		c.Locals("user_id", userID)
		c.Locals("user_email", claims["email"])
		c.Locals("user_role", claims["role"])

		// Izinkan request lanjut ke Handler berikutnya
		return c.Next()
	}
}
