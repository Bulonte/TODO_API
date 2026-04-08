package token

import (
	"TODO_API/pkg/cache"
	"context"
	"fmt"
	"time"
)

// TokenManager token 管理器接口
type TokenManager interface {
	// AddToBlacklist 将 token 添加到黑名单
	AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error
	// IsBlacklisted 检查 token 是否在黑名单中
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

// RedisTokenManager Redis token 管理器实现
type RedisTokenManager struct {
	redisClient cache.RedisClient
}

// NewRedisTokenManager 创建 Redis token 管理器实例
func NewRedisTokenManager(redisClient cache.RedisClient) TokenManager {
	return &RedisTokenManager{
		redisClient: redisClient,
	}
}

// getBlacklistKey 生成黑名单键
func getBlacklistKey(token string) string {
	return fmt.Sprintf("token:blacklist:%s", token)
}

// AddToBlacklist 将 token 添加到黑名单
func (m *RedisTokenManager) AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	// 生成黑名单键
	blacklistKey := getBlacklistKey(token)
	
	// 设置 token 到黑名单
	return m.redisClient.Set(ctx, blacklistKey, "1", expiration)
}

// IsBlacklisted 检查 token 是否在黑名单中
func (m *RedisTokenManager) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	// 生成黑名单键
	blacklistKey := getBlacklistKey(token)
	
	// 检查 token 是否在黑名单中
	exists, err := m.redisClient.Exists(ctx, blacklistKey)
	if err != nil {
		return false, err
	}
	
	return exists > 0, nil
}
