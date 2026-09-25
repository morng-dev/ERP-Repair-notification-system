package entities

import "github.com/google/uuid"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Firsname string `json:"first_name" validate:"required"`
	Lastname string `json:"last_name" validate:"required"`
	Avatar   string `json:"avatar,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user,omitempty"`
}

type ChangePasswordRequest struct {
	Oldpassword string `json:"old_password"`
	Newpassword string `json:"new_password"`
}

type RefreshTokenRequest struct {
	Refreshtoken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}

type UserResponse struct {
	ID        uuid.UUID     `json:"id"`
	Email     string        `json:"email"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Avatar    string        `json:"avatar"`
	Active    bool          `json:"active"`
	Role      *RoleResponse `json:"role,omitempty"`
}

type RoleResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
