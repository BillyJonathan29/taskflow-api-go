package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// DTO untuk payload saat user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
