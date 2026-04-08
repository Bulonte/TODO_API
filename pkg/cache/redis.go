package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis 客户端接口
type RedisClient interface {
	// 设置键值对
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// 获取值
	Get(ctx context.Context, key string) (string, error)
	// 删除键
	Del(ctx context.Context, keys ...string) error
	// 检查键是否存在
	Exists(ctx context.Context, keys ...string) (int64, error)
	// 健康检查
	Ping(ctx context.Context) error
	// 关闭连接
	Close() error
}

// redisClient Redis 客户端实现
type redisClient struct {
	client *redis.Client
}

// NewRedisClient 创建 Redis 客户端实例
func NewRedisClient(addr, password string, db int) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10,
	})

	return &redisClient{client: client}
}

// Set 设置键值对
func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del 删除键
func (r *redisClient) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *redisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Ping 健康检查
func (r *redisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close 关闭连接
func (r *redisClient) Close() error {
	return r.client.Close()
}
