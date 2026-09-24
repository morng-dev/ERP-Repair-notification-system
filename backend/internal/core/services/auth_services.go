package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/cache"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"github.com/morng-dev/erp/pkg/utils"
)

type AuthService struct {
	userRepo      repositories.UserRepository
	roleRepo      repositories.RoleRepository
	UserRedisRepo cache.UserRedisRepo
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, UserRedisRepo cache.UserRedisRepo) services.AuthService {
	return &AuthService{
		userRepo:      userRepo,
		roleRepo:      roleRepo,
		UserRedisRepo: UserRedisRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req *entities.RegisterRequest) (*entities.User, error) {
	exists, err := s.userRepo.GetEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	userRole, err := s.roleRepo.GetByName(ctx, "user")
	if err != nil {
		return nil, errors.New("ไม่พบบทบาทผู้ใช้")
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Email:    req.Email,
		Firsname: req.Firsname,
		Lastname: req.Lastname,
		Avatar:   req.Avatar,
		RoleID:   userRole.ID,
	}

	result, err := s.userRepo.Create(ctx, user, hashPassword)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *AuthService) Login(ctx context.Context, req *entities.LoginRequest) (*entities.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("ผู้ใช้หรือรหัสผ่านไม่ถูกต้อง")
	}
	if !user.Active {
		return nil, errors.New("ถูกระงับผู้ใช้โปรดติดต่อผู้ให้บริการ")
	}
	hashPassword, err := s.userRepo.GetPasswordHash(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	if !utils.CompairPassword(req.Password, hashPassword) {
		return nil, errors.New("ผู้ใช้หรือรหัสผ่านไม่ถูกต้อง")
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Email, user.Role.Name)
	if err != nil {
		return nil, err
	}

	return &entities.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *entities.ChangePasswordRequest) error {
	hashPassword, err := s.userRepo.GetPasswordHash(ctx, userID)
	if err != nil {
		return err
	}
	if !utils.CompairPassword(req.Oldpassword, hashPassword) {
		return errors.New("password invalid")
	}

	newPasswordHash, err := utils.HashPassword(req.Newpassword)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePassword(ctx, userID, newPasswordHash)
}

func (s *AuthService) RefreshToken(ctx context.Context, req *entities.RefreshTokenRequest) (*entities.LoginResponse, error) {
	user, err := s.userRepo.GetByRefreshToken(ctx, req.Refreshtoken)
	if err != nil {
		return nil, err
	}
	token, err := utils.GenerateToken(user.Email, user.ID.String(), user.Role.ID.String())
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.GenerrateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.SetRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, err
	}
	return &entities.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GenerrateRefreshToken() (string, error) {
	byte := make([]byte, 16)
	if _, err := rand.Read(byte); err != nil {
		return "", err
	}
	return hex.EncodeToString(byte), nil
}
