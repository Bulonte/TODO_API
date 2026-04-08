package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// AppError 应用错误类型
type AppError struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
	Err     error  `json:"err"`     // 原始错误
	Stack   string `json:"stack"`   // 堆栈信息
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新的应用错误
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Stack:   getStackTrace(),
	}
}

// BadRequest 创建 400 错误
func BadRequest(message string, err error) *AppError {
	return New(400, message, err)
}

// Unauthorized 创建 401 错误
func Unauthorized(message string, err error) *AppError {
	return New(401, message, err)
}

// Forbidden 创建 403 错误
func Forbidden(message string, err error) *AppError {
	return New(403, message, err)
}

// NotFound 创建 404 错误
func NotFound(message string, err error) *AppError {
	return New(404, message, err)
}

// InternalServerError 创建 500 错误
func InternalServerError(message string, err error) *AppError {
	return New(500, message, err)
}

// getStackTrace 获取堆栈信息
func getStackTrace() string {
	var stack []string
	for i := 2; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = append(stack, fmt.Sprintf("%s:%d", file, line))
	}
	return strings.Join(stack, "\n")
}
