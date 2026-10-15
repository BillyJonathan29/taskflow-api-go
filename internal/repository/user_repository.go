package repository

import (
	"taskflow-api/internal/domain"

	"gorm.io/gorm"
)

// Kontrak kerja (Interface): Memberitahu fungsi apa saja yang harus dimiliki UserRepository
type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

// Struct implementasi yang memegang koneksi GORM
type userRepository struct {
	db *gorm.DB
}

// Constructor: fungsi untuk membuat instance userRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
