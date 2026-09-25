package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type AuthService interface {
	Register(ctx context.Context, req *entities.RegisterRequest) (*entities.User, error)
	Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *entities.ChangePasswordRequest) error
	RefreshToken(ctx context.Context, req *entities.RefreshTokenRequest) (*entities.LoginResponse, error)
	ForgotPassword(ctx context.Context, req *entities.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req *entities.ResetPasswordRequest) error
}
