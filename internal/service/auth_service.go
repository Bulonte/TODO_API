package service

import (
	"TODO_API/config"
	"TODO_API/internal/app/dto/request"
	"TODO_API/internal/app/dto/response"
	"TODO_API/internal/domain/model"
	"TODO_API/internal/repository"
	"TODO_API/pkg/encryption"
	"TODO_API/pkg/jwt"
	"context"
	"errors"
	"time"
)

// 认证服务接口
type AuthService interface {
	Register(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, error)
	Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*response.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// 生成携带Token的AuthResponse用户信息
func (s *authService) generateAuthServiceWithToken(user *model.User) (*response.AuthResponse, error) {
	//生成访问令牌
	accessToken, err := jwt.GenerateToken(user.ID, user.Username, false)
	if err != nil {
		return nil, err
	}

	//生成刷新令牌
	refreshToken, err := jwt.GenerateToken(user.ID, user.Username, true)
	if err != nil {
		return nil, err
	}

	//计算过期时间
	expiresAt := time.Now().Add(time.Duration(config.GlobalConfig.JWT.AccessExpire) * time.Second).Unix()

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Avatar:   user.AvatarURL,
		},
	}, nil
}

// 用户注册
func (s *authService) Register(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, error) {
	//检查用户名是否存在
	exitingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if exitingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	//检查邮箱是否存在
	exitingUser, _ = s.userRepo.GetByEmail(ctx, req.Email)
	if exitingUser != nil {
		return nil, errors.New("邮箱已存在")
	}

	//密码加密
	HashedPassword, err := encryption.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	//创建用户
	User := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: HashedPassword,
		Status:       1,
	}

	//err = s.userRepo.Create(ctx, User)
	if err = s.userRepo.Create(ctx, User); err != nil {
		return nil, err
	}

	//返回带token的用户信息
	return s.generateAuthServiceWithToken(User)
}

// 用户登录
func (s *authService) Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error) {
	//获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	//验证密码
	if !encryption.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("用户名或密码错误")
	}

	//检查用户状态
	if user.Status == 0 {
		return nil, errors.New("用户已被封禁")
	}

	//返回带token的用户信息
	return s.generateAuthServiceWithToken(user)
}

// 刷新令牌方法
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*response.AuthResponse, error) {
	newAccessToken, err := jwt.RefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("刷新令牌无效")
	}

	//解析令牌
	claims, err := jwt.ParseToken(newAccessToken)
	if err != nil {
		return nil, err
	}

	//获取用户信息
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}

	expiresAt := time.Now().Add(time.Duration(config.GlobalConfig.JWT.AccessExpire) * time.Second).Unix()

	return &response.AuthResponse{
		AccessToken: newAccessToken,
		ExpiresAt:   expiresAt,
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Avatar:   user.AvatarURL,
		},
	}, nil
}
