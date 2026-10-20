package auth

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// Kontrak interface TokenManager
type TokenManager interface {
	GenerateToken(userID uint, email, role string) (string, error)
	StoreSession(ctx context.Context, userID uint, token string) error
	ValidateSession(ctx context.Context, userID uint, token string) bool
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type tokenManager struct {
	secretKey []byte
	expHours  time.Duration
	rdb       *redis.Client
}

// Constructor untuk membuat TokenManager
func NewTokenManager(rdb *redis.Client) TokenManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_key"
	}
	expHoursStr := os.Getenv("JWT_EXP_HOURS")
	expHours, err := strconv.Atoi(expHoursStr)
	if err != nil {
		expHours = 24 // Fallback jika tidak diset, default 24 jam
	}

	return &tokenManager{
		secretKey: []byte(secret),
		expHours:  time.Duration(expHours) * time.Hour,
		rdb:       rdb,
	}
}

// 1. Fungsi untuk mencetak token JWT yang berisi payload (Claims)
func (m *tokenManager) GenerateToken(userID uint, email, role string) (string, error) {
	// Claims adalah payload data yang disisipkan di dalam token
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(m.expHours).Unix(), // Waktu kedaluwarsa token (Unix timestamp)
		"iat":     time.Now().Unix(),                 // Waktu token diterbitkan (Issued At)
	}

	// Buat token dengan algoritma enkripsi HMAC SHA-256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Tanda tangani token menggunakan secret key kita
	return token.SignedString(m.secretKey)

}

// 2. Fungsi untuk menyimpan sesi login token ke memori RAM Redis
func (m *tokenManager) StoreSession(ctx context.Context, userID uint, token string) error {
	// Key di redis, misal: "session:user:1"
	redisKey := fmt.Sprintf("Session:user:%d", userID)
	// Set data ke Redis dengan TTL (Time To Live) otomatis terhapus saat expired
	return m.rdb.Set(ctx, redisKey, token, m.expHours).Err()
}

// 3. Fungsi untuk memeriksa apakah token user masih aktif di Redis
func (m *tokenManager) ValidateSession(ctx context.Context, userID uint, token string) bool {
	redisKey := fmt.Sprintf("Session:user:%d", userID)
	// Ambil token yang tersimpan di Redis
	savedToken, err := m.rdb.Get(ctx, redisKey).Result()

	if err != nil {
		return false // Token tidak ada di Redis (sudah expired atau di-logout)
	}
	// Cocokkan apakah token yang dibawa user sama dengan token yang tersimpan
	return savedToken == token
}

// 4. Fungsi untuk membedah dan memverifikasi keaslian signature JWT
func (m *tokenManager) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	// Parse token dan verifikasi keaslian tandatangan (signature)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Pastikan algoritma yang dipakai tetap HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Algoritma signing tidak valid: %v", t.Header["alg"])
		}
		return m.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Ekstrak data claims jika token valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("Token tidak valid")
}
