package service

import (
	"TODO_API/internal/app/dto/request"
	"TODO_API/internal/app/dto/response"
	"TODO_API/internal/repository"
	"TODO_API/pkg/encryption"
	"context"
	"errors"
)

type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*response.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *request.UpdateProfileRequest) (*response.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *userService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}

	return &response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.AvatarURL,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint, req *request.UpdateProfileRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}

	if req.Email != "" {
		//检查邮箱是否被使用
		existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
		if existingUser != nil && existingUser.ID != user.ID {
			return nil, errors.New("邮箱已被使用")
		}
		user.Email = req.Email
	}
	if req.AvatarURL != "" {
		user.AvatarURL = &req.AvatarURL
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return &response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.AvatarURL,
	}, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}
	//验证旧密码
	if !encryption.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("旧密码错误")
	}

	//加密新密码
	hashedPassword, err := encryption.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword
	return s.userRepo.Update(ctx, user)
}
