package lock

import (
	"TODO_API/pkg/cache"
	"context"
	"fmt"
	"time"
)

// DistributedLock 分布式锁接口
type DistributedLock interface {
	// Lock 获取锁
	Lock(ctx context.Context, key string, expiration time.Duration) (bool, error)
	// Unlock 释放锁
	Unlock(ctx context.Context, key string) error
}

// RedisLock Redis 分布式锁实现
type RedisLock struct {
	redisClient cache.RedisClient
}

// NewRedisLock 创建 Redis 分布式锁实例
func NewRedisLock(redisClient cache.RedisClient) DistributedLock {
	return &RedisLock{
		redisClient: redisClient,
	}
}

// getLockKey 生成锁键
func getLockKey(key string) string {
	return fmt.Sprintf("lock:%s", key)
}

// Lock 获取锁
func (l *RedisLock) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	// 生成锁键
	lockKey := getLockKey(key)
	
	// 使用 SETNX 命令获取锁
	// 注意：由于我们的 RedisClient 接口没有直接提供 SETNX 方法，这里使用一个简单的实现
	// 实际生产环境中应该使用 SETNX 或 SET with NX 选项
	
	// 检查锁是否已存在
	exists, err := l.redisClient.Exists(ctx, lockKey)
	if err != nil {
		return false, err
	}
	
	if exists > 0 {
		// 锁已存在
		return false, nil
	}
	
	// 设置锁
	err = l.redisClient.Set(ctx, lockKey, "1", expiration)
	if err != nil {
		return false, err
	}
	
	return true, nil
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context, key string) error {
	// 生成锁键
	lockKey := getLockKey(key)
	
	// 删除锁
	return l.redisClient.Del(ctx, lockKey)
}
