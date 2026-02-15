package request

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Email     string `json:"email,omitempty" binding:"email"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password,omitempty" binding:"required"`
	NewPassword     string `json:"new_password,omitempty" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
