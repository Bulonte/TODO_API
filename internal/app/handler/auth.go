package handler

import (
	"TODO_API/internal/app/dto/request"
	"TODO_API/internal/service"
	"TODO_API/pkg/response"

	_ "TODO_API/docs" // 导入Swagger文档

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

// 创建AuthHandler实例
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=response.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	//绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	//调用服务层进行注册
	authResp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "用户名已存在", "邮箱已存在":
			response.BadRequest(c, err.Error())
		default:
			response.InternalServerError(c, "注册失败"+err.Error())
		}
		return
	}

	response.Success(c, authResp)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=response.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	//绑定参数并验证
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	//调用服务层进行登录
	authResp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "用户名或密码错误" {
			response.Unauthorized(c, err.Error())
		} else {
			response.InternalServerError(c, "登录失败"+err.Error())
		}
		return
	}

	response.Success(c, authResp)
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌和刷新令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.RefreshTokenRequest true "刷新令牌信息"
// @Success 200 {object} response.Response{data=response.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误"+err.Error())
		return
	}

	//调用服务层刷新令牌
	authResp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "令牌刷新失败"+err.Error())
		return
	}

	response.Success(c, authResp)
}
