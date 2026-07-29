package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/repositories"
	"github.com/morng-dev/erp/internal/core/domain/ports/services"
	"github.com/morng-dev/erp/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
	rdb      *redis.Client
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, rdb *redis.Client) services.AuthService {
	return &AuthService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		rdb:      rdb,
	}
}

func (s *AuthService) Register(ctx context.Context, req *entities.RegisterRequest) (*entities.User, error) {
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("มีผู้ใช้แล้วในระบบ")
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

	return s.userRepo.GetByID(ctx, user.ID)
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
		User:  *user,
	}, nil
}

func (s *AuthService) CachePermissions(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	key := "user:permissions:" + userID.String()
	if err := s.rdb.Del(ctx, key); err != nil {
		log.Printf("cannot clear cache:%v", err)
	}
	pers := make([]interface{}, 0, len(user.Role.Permissions))
	for _, p := range user.Role.Permissions {
		pers = append(pers, p.Name)
	}
	if len(pers) > 0 {
		if err := s.rdb.SAdd(ctx, key, pers...).Err(); err != nil {
			return err
		}
		if err := s.rdb.Expire(ctx, key, 1*time.Hour).Err(); err != nil {
			return err
		}
	}
	return nil
}
