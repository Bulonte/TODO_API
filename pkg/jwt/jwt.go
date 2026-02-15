package jwt

import (
	"TODO_API/config"
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

// 初始化JWT配置
func InitJWT() {
	if config.GlobalConfig.JWT.Secret == "" {
		panic("JWT密钥未配置")
	}
	jwtSecret = []byte(config.GlobalConfig.JWT.Secret)
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

// 刷新令牌
func RefreshToken(tokenString string) (string, error) {
	refreshToken, err := ParseToken(tokenString)
	if err != nil {
		return "", errors.New("刷新令牌无效")
	}

	return GenerateToken(refreshToken.UserID, refreshToken.Username, false)
}
