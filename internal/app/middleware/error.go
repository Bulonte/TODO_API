package middleware

import (
	"TODO_API/pkg/errors"
	"TODO_API/pkg/logger"
	"TODO_API/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// 检查是否是 AppError 类型
				if appErr, ok := e.Err.(*errors.AppError); ok {
					// 记录错误日志
					logger.Error("应用错误",
						zap.Int("code", appErr.Code),
						zap.String("message", appErr.Message),
						zap.String("stack", appErr.Stack),
					)

					// 返回错误响应
					response.Error(c, appErr.Code, appErr.Message)
					return
				} else {
					// 处理普通错误
					logger.Error("未知错误", zap.Error(e.Err))
					response.InternalServerError(c, "服务器内部错误")
					return
				}
			}
		}
	}
}
