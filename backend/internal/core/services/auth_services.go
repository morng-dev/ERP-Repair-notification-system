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
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo      repositories.UserRepository
	roleRepo      repositories.RoleRepository
	mailer        services.Mailer
	UserRedisRepo cache.UserRedisRepo
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, UserRedisRepo cache.UserRedisRepo, mailer services.Mailer) services.AuthService {
	return &AuthService{
		userRepo:      userRepo,
		roleRepo:      roleRepo,
		mailer:        mailer,
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

	if err := s.userRepo.Create(ctx, user, hashPassword); err != nil {
		return nil, err
	}
	user, err = s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
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

	token, err := utils.GenerateToken(user.Email, user.ID.String(), user.Role.Name)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.GenerateToken()
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.SetRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, err
	}

	return &entities.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: &entities.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.Firsname,
			LastName:  user.Lastname,
			Avatar:    user.Avatar,
			Active:    user.Active,
			Role: &entities.RoleResponse{
				Name:        user.Role.Name,
				Description: user.Role.Description,
			},
		},
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *entities.ChangePasswordRequest) error {
	passwordHash, err := s.userRepo.GetPasswordHash(ctx, userID)
	if err != nil {
		return err
	}
	if !utils.CompairPassword(req.Oldpassword, passwordHash) {
		return errors.New("password invalid")
	}
	newPasswordhash, err := utils.HashPassword(req.Newpassword)
	if err != nil {
		return err
	}
	if err := s.userRepo.UpdatePassword(ctx, userID, newPasswordhash); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req *entities.ForgotPasswordRequest) error {
	exists, err := s.userRepo.GetEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	resetToken, err := s.GenerateToken()
	if err != nil {
		return err
	}
	if err := s.userRepo.SetResetToken(ctx, req.Email, resetToken); err != nil {
		return err
	}
	if err := s.mailer.SendResetPassword(ctx, req.Email, resetToken); err != nil {
		return errors.New("can'n send email plese retry again")
	}
	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *entities.ResetPasswordRequest) error {
	user, err := s.userRepo.GetByResetToken(ctx, req.ResetToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Invalid or expiry token")
		}
		return err
	}
	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, passwordHash); err != nil {
		return err
	}
	if err := s.userRepo.ClearResetToken(ctx, user.ID); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *entities.RefreshTokenRequest) (*entities.LoginResponse, error) {
	user, err := s.userRepo.GetByRefreshToken(ctx, req.Refreshtoken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token invalid")
		}
		return nil, err
	}
	refreshToken, err := s.GenerateToken()
	if err != nil {
		return nil, err
	}
	token, err := utils.GenerateToken(user.Email, user.ID.String(), user.Role.Name)
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

func (s *AuthService) GenerateToken() (string, error) {
	byte := make([]byte, 32)
	if _, err := rand.Read(byte); err != nil {
		return "", err
	}
	return hex.EncodeToString(byte), nil
}
