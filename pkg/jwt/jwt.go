package jwt

import (
	"TODO_API/config"
	"TODO_API/pkg/token"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret []byte
var tokenManager token.TokenManager

// 初始化JWT配置
func InitJWT() {
	if config.GlobalConfig.JWT.Secret == "" {
		panic("JWT密钥未配置")
	}
	jwtSecret = []byte(config.GlobalConfig.JWT.Secret)
}

// SetTokenManager 设置token管理器
func SetTokenManager(manager token.TokenManager) {
	tokenManager = manager
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string, isRefresh bool) (string, error) {
	nowTime := time.Now()
	var expireTime time.Time

	if isRefresh {
		// 刷新令牌有效期更长
		expireTime = nowTime.Add(time.Duration(config.GlobalConfig.JWT.RefreshExpire) * time.Second)
	} else {
		// 访问令牌
		expireTime = nowTime.Add(time.Duration(config.GlobalConfig.JWT.AccessExpire) * time.Second)
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.GlobalConfig.JWT.Issuer,
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 检查token是否在黑名单中
	if tokenManager != nil {
		isBlacklisted, err := tokenManager.IsBlacklisted(context.Background(), tokenString)
		if err != nil {
			return nil, err
		}
		if isBlacklisted {
			return nil, errors.New("令牌已被注销")
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// RefreshToken 刷新令牌
func RefreshToken(tokenString string) (string, error) {
	refreshToken, err := ParseToken(tokenString)
	if err != nil {
		return "", errors.New("刷新令牌无效")
	}

	// 将旧的刷新令牌添加到黑名单
	if tokenManager != nil {
		expiration := time.Duration(config.GlobalConfig.JWT.RefreshExpire) * time.Second
		tokenManager.AddToBlacklist(context.Background(), tokenString, expiration)
	}

	return GenerateToken(refreshToken.UserID, refreshToken.Username, false)
}

// RevokeToken 注销令牌
func RevokeToken(tokenString string) error {
	if tokenManager == nil {
		return errors.New("token管理器未初始化")
	}

	// 解析令牌获取过期时间
	claims, err := ParseToken(tokenString)
	if err != nil {
		return err
	}

	// 计算剩余过期时间
	expiration := time.Until(claims.ExpiresAt.Time)
	if expiration <= 0 {
		expiration = time.Second // 至少设置1秒过期
	}

	// 将令牌添加到黑名单
	return tokenManager.AddToBlacklist(context.Background(), tokenString, expiration)
}
