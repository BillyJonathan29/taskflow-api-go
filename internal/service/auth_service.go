package service

import (
	"context"
	"errors"

	"taskflow-api/internal/auth"
	"taskflow-api/internal/dto"
	"taskflow-api/internal/model"
	"taskflow-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// kontrak interface: menambah method login
type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo     repository.UserRepository
	tokenManager auth.TokenManager // Dependency ke TokenManager

}

// Constructor AuthService (Menerima TOkenManager)
func NewAuthService(userRepo repository.UserRepository, tokenManager auth.TokenManager) AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}

}

func (s *authService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// 1. Cek apakah email sudah terdaftar
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("Email sudah terdaftar")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 2. Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "member",
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

// Method Login: Mencocokkan password, buat JWT, simpan ke Redis
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. Cari user berdasarkan email di database
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// Demi keamanan, jika email tidak ada, pesan error disamarkan
		return nil, errors.New("Email atau password salah")
	}

	// 2. Cocokkan hash password di DB dengan plaintext password yang diketik user
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("Email atau password salah")
	}

	// 3. Cetak token JWT
	token, err := s.tokenManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("Gagal memproses token autentikasi")
	}

	// 4. Simpan sesi login ke Redis (Skalabilitas & Manajemen Session)
	if err := s.tokenManager.StoreSession(ctx, user.ID, token); err != nil {
		return nil, errors.New("Gagal menyimpan sesi ke redis")
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil

}
