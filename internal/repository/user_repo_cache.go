package repository

import (
	"TODO_API/internal/domain/model"
	"TODO_API/pkg/cache"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// UserRepositoryCache 用户仓库缓存装饰器
type UserRepositoryCache struct {
	userRepo   UserRepository
	redisCache cache.RedisClient
}

// NewUserRepositoryCache 创建带缓存的用户仓库
func NewUserRepositoryCache(userRepo UserRepository, redisCache cache.RedisClient) UserRepository {
	return &UserRepositoryCache{
		userRepo:   userRepo,
		redisCache: redisCache,
	}
}

// getUserCacheKey 生成用户缓存键
func getUserCacheKey(id uint) string {
	return fmt.Sprintf("user:id:%d", id)
}

// getUsernameCacheKey 生成用户名缓存键
func getUsernameCacheKey(username string) string {
	return fmt.Sprintf("user:username:%s", username)
}

// getEmailCacheKey 生成邮箱缓存键
func getEmailCacheKey(email string) string {
	return fmt.Sprintf("user:email:%s", email)
}

// Create 创建用户
func (r *UserRepositoryCache) Create(ctx context.Context, user *model.User) error {
	err := r.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}
	
	// 缓存用户信息（包含 PasswordHash）
	cacheData := map[string]interface{}{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"password_hash": user.PasswordHash,
		"status":       user.Status,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	}
	
	userJSON, err := json.Marshal(cacheData)
	if err == nil {
		r.redisCache.Set(ctx, getUserCacheKey(user.ID), userJSON, 24*time.Hour)
		r.redisCache.Set(ctx, getUsernameCacheKey(user.Username), userJSON, 24*time.Hour)
		r.redisCache.Set(ctx, getEmailCacheKey(user.Email), userJSON, 24*time.Hour)
	}
	
	return nil
}

// GetByID 通过ID获取用户
func (r *UserRepositoryCache) GetByID(ctx context.Context, id uint) (*model.User, error) {
	// 尝试从缓存获取
	cacheKey := getUserCacheKey(id)
	userJSON, err := r.redisCache.Get(ctx, cacheKey)
	if err == nil {
		var cacheData map[string]interface{}
		if json.Unmarshal([]byte(userJSON), &cacheData) == nil {
			user := &model.User{
				ID:           uint(cacheData["id"].(float64)),
				Username:     cacheData["username"].(string),
				Email:        cacheData["email"].(string),
				PasswordHash: cacheData["password_hash"].(string),
				Status:       uint8(cacheData["status"].(float64)),
			}
			// 解析日期时间字段
			if createdAtStr, ok := cacheData["created_at"].(string); ok {
				if createdAt, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
					user.CreatedAt = createdAt
				}
			}
			if updatedAtStr, ok := cacheData["updated_at"].(string); ok {
				if updatedAt, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
					user.UpdatedAt = updatedAt
				}
			}
			return user, nil
		}
	}
	
	// 缓存未命中，从数据库获取
	user, err := r.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// 缓存用户信息（包含 PasswordHash）
	if user != nil {
		cacheData := map[string]interface{}{
			"id":           user.ID,
			"username":     user.Username,
			"email":        user.Email,
			"password_hash": user.PasswordHash,
			"status":       user.Status,
			"created_at":   user.CreatedAt,
			"updated_at":   user.UpdatedAt,
		}
		
		userJSON, err := json.Marshal(cacheData)
		if err == nil {
			r.redisCache.Set(ctx, cacheKey, userJSON, 24*time.Hour)
		}
	}
	
	return user, nil
}

// GetByUsername 通过用户名获取用户
func (r *UserRepositoryCache) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	// 尝试从缓存获取
	cacheKey := getUsernameCacheKey(username)
	userJSON, err := r.redisCache.Get(ctx, cacheKey)
	if err == nil {
		var cacheData map[string]interface{}
		if json.Unmarshal([]byte(userJSON), &cacheData) == nil {
			user := &model.User{
				ID:           uint(cacheData["id"].(float64)),
				Username:     cacheData["username"].(string),
				Email:        cacheData["email"].(string),
				PasswordHash: cacheData["password_hash"].(string),
				Status:       uint8(cacheData["status"].(float64)),
			}
			// 解析日期时间字段
			if createdAtStr, ok := cacheData["created_at"].(string); ok {
				if createdAt, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
					user.CreatedAt = createdAt
				}
			}
			if updatedAtStr, ok := cacheData["updated_at"].(string); ok {
				if updatedAt, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
					user.UpdatedAt = updatedAt
				}
			}
			return user, nil
		}
	}
	
	// 缓存未命中，从数据库获取
	user, err := r.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	
	// 缓存用户信息（包含 PasswordHash）
	if user != nil {
		cacheData := map[string]interface{}{
			"id":           user.ID,
			"username":     user.Username,
			"email":        user.Email,
			"password_hash": user.PasswordHash,
			"status":       user.Status,
			"created_at":   user.CreatedAt,
			"updated_at":   user.UpdatedAt,
		}
		
		userJSON, err := json.Marshal(cacheData)
		if err == nil {
			r.redisCache.Set(ctx, cacheKey, userJSON, 24*time.Hour)
			r.redisCache.Set(ctx, getUserCacheKey(user.ID), userJSON, 24*time.Hour)
			r.redisCache.Set(ctx, getEmailCacheKey(user.Email), userJSON, 24*time.Hour)
		}
	}
	
	return user, nil
}

// GetByEmail 通过邮箱获取用户
func (r *UserRepositoryCache) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// 尝试从缓存获取
	cacheKey := getEmailCacheKey(email)
	userJSON, err := r.redisCache.Get(ctx, cacheKey)
	if err == nil {
		var cacheData map[string]interface{}
		if json.Unmarshal([]byte(userJSON), &cacheData) == nil {
			user := &model.User{
				ID:           uint(cacheData["id"].(float64)),
				Username:     cacheData["username"].(string),
				Email:        cacheData["email"].(string),
				PasswordHash: cacheData["password_hash"].(string),
				Status:       uint8(cacheData["status"].(float64)),
			}
			// 解析日期时间字段
			if createdAtStr, ok := cacheData["created_at"].(string); ok {
				if createdAt, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
					user.CreatedAt = createdAt
				}
			}
			if updatedAtStr, ok := cacheData["updated_at"].(string); ok {
				if updatedAt, err := time.Parse(time.RFC3339, updatedAtStr); err == nil {
					user.UpdatedAt = updatedAt
				}
			}
			return user, nil
		}
	}
	
	// 缓存未命中，从数据库获取
	user, err := r.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	
	// 缓存用户信息（包含 PasswordHash）
	if user != nil {
		cacheData := map[string]interface{}{
			"id":           user.ID,
			"username":     user.Username,
			"email":        user.Email,
			"password_hash": user.PasswordHash,
			"status":       user.Status,
			"created_at":   user.CreatedAt,
			"updated_at":   user.UpdatedAt,
		}
		
		userJSON, err := json.Marshal(cacheData)
		if err == nil {
			r.redisCache.Set(ctx, cacheKey, userJSON, 24*time.Hour)
			r.redisCache.Set(ctx, getUserCacheKey(user.ID), userJSON, 24*time.Hour)
			r.redisCache.Set(ctx, getUsernameCacheKey(user.Username), userJSON, 24*time.Hour)
		}
	}
	
	return user, nil
}

// Update 更新用户
func (r *UserRepositoryCache) Update(ctx context.Context, user *model.User) error {
	err := r.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}
	
	// 更新缓存（包含 PasswordHash）
	cacheData := map[string]interface{}{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"password_hash": user.PasswordHash,
		"status":       user.Status,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	}
	
	userJSON, err := json.Marshal(cacheData)
	if err == nil {
		r.redisCache.Set(ctx, getUserCacheKey(user.ID), userJSON, 24*time.Hour)
		r.redisCache.Set(ctx, getUsernameCacheKey(user.Username), userJSON, 24*time.Hour)
		r.redisCache.Set(ctx, getEmailCacheKey(user.Email), userJSON, 24*time.Hour)
	}
	
	return nil
}
