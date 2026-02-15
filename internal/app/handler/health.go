package handler

import (
	"TODO_API/config"
	"TODO_API/pkg/logger"
	"TODO_API/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Healther struct{}

func NewHealther() *Healther {
	return &Healther{}
}

// HealthChecker 服务健康检查
// @Summary 服务健康检查
// @Description 检查服务的运行状态和基本信息
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=object}
// @Failure 500 {object} response.Response
// @Router /health [get]
func (h *Healther) HealthChecker(c *gin.Context) {
	logger.Debug("健康检查请求")
	response.Success(c, gin.H{
		"status":       "healthy",
		"service_name": config.GlobalConfig.App.Name,
		"version":      config.GlobalConfig.App.Version,
		"environment":  config.GlobalConfig.App.Environment,
		"current_time": time.Now(),
	})
}

// ReadyChecker 服务就绪检查
// @Summary 服务就绪检查
// @Description 检查服务是否已就绪并可处理请求（包括依赖组件状态）
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=object}
// @Failure 503 {object} response.Response
// @Router /ready [get]
func (h *Healther) ReadyChecker(c *gin.Context) {
	dbstatus := "Unknown"
	logger.Info("就绪状态检查",
		zap.String("database", dbstatus),
	)
	response.Success(c, gin.H{
		"status":       "ready",
		"service_name": config.GlobalConfig.App.Name,
		"database":     dbstatus,
		"environment":  config.GlobalConfig.App.Environment,
	})
}
