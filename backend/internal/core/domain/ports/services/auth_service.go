package services

import (
	"context"

	"github.com/morng-dev/erp/internal/core/domain/entities"
)

type AuthService interface {
	Register(ctx context.Context, req *entities.RegisterRequest) (*entities.User, error)
	Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error)
}
