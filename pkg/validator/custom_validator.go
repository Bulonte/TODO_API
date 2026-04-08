package validator

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitValidator 初始化自定义验证器
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册自定义验证器
		v.RegisterValidation("password", validatePassword)
		v.RegisterValidation("future_date", validateFutureDate)
		v.RegisterValidation("date_range", validateDateRange)
	}
}

// validatePassword 验证密码强度
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	// 检查是否包含至少一个数字
	hasDigit := false
	// 检查是否包含至少一个字母
	hasLetter := false
	// 检查是否包含至少一个特殊字符
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= '0' && char <= '9':
			hasDigit = true
		case (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z'):
			hasLetter = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}

	return hasDigit && hasLetter && hasSpecial
}

// validateFutureDate 验证日期是否为未来日期
func validateFutureDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return date.After(time.Now())
}

// validateDateRange 验证日期范围
func validateDateRange(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// 检查日期是否在合理范围内（例如：不超过10年）
	maxDate := time.Now().AddDate(10, 0, 0)
	return date.Before(maxDate)
}

// GetValidationError 获取验证错误信息
func GetValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s 是必填字段", e.Field())
			case "min":
				return fmt.Sprintf("%s 长度不能小于 %s", e.Field(), e.Param())
			case "max":
				return fmt.Sprintf("%s 长度不能大于 %s", e.Field(), e.Param())
			case "email":
				return fmt.Sprintf("%s 不是有效的邮箱地址", e.Field())
			case "oneof":
				return fmt.Sprintf("%s 必须是以下值之一: %s", e.Field(), e.Param())
			case "password":
				return "密码必须至少包含8个字符，且包含至少一个数字和一个字母"
			case "future_date":
				return "日期必须是未来日期"
			case "date_range":
				return "日期超出合理范围"
			default:
				return fmt.Sprintf("%s 验证失败: %s", e.Field(), e.Tag())
			}
		}
	}
	return "输入验证失败"
}
