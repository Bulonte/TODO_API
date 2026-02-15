package handler

import (
	"TODO_API/internal/app/dto/request"
	"TODO_API/internal/app/middleware"
	"TODO_API/internal/service"
	"TODO_API/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandeler struct {
	userService service.UserService
}

func NewUserHandeler(userService service.UserService) *UserHandeler {
	return &UserHandeler{userService: userService}
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前已登录用户的个人信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=response.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/profile [get]
func (u *UserHandeler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)

	user, err := u.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "获取用户信息失败")
		return
	}
	response.Success(c, user)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前用户的个人信息（如用户名、邮箱等）
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.UpdateProfileRequest true "用户信息更新请求"
// @Success 200 {object} response.Response{data=response.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/profile [put]
func (u *UserHandeler) UpdateProfile(c *gin.Context) {
	var req request.UpdateProfileRequest
	userID := middleware.GetUserIDFromContext(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	user, err := u.userService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		response.InternalServerError(c, "修改用户信息失败")
		return
	}
	response.Success(c, user)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的登录密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.ChangePasswordRequest true "密码修改请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /user/change-password [put]
func (u *UserHandeler) ChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest
	userID := middleware.GetUserIDFromContext(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	err := u.userService.ChangePassword(c.Request.Context(), userID, &req)
	if err != nil {
		response.InternalServerError(c, "修改密码失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
