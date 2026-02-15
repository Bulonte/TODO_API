package middleware

import (
	"TODO_API/pkg/jwt"
	"TODO_API/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			response.Unauthorized(c, "请提供认证令牌")
			c.Abort()
			return
		}

		//检查token格式 Bearer <token>
		parts := strings.SplitN(authToken, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "令牌格式错误，应为: Bearer <token>")
			c.Abort()
			return
		}

		// 解析token
		tokenString := parts[1]
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "令牌无效或已过期"+err.Error())
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("UserID", claims.UserID)
		c.Set("UserName", claims.Username)
		// 记录日志
		zap.L().Debug("用户认证成功",
			zap.Uint("UserID", claims.UserID),
			zap.String("Username", claims.Username),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	}
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) uint {
	if userID, exists := c.Get("UserID"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

// GetUserNameFromContext 从上下文中获取用户信息
func GetUserNameFromContext(c *gin.Context) string {
	if userName, exists := c.Get("UserName"); exists {
		if name, ok := userName.(string); ok {
			return name
		}
	}
	return ""
}
