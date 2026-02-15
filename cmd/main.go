package main

import (
	"TODO_API/config"
	"TODO_API/internal/app/handler"
	"TODO_API/internal/app/middleware"
	"TODO_API/internal/repository"
	"TODO_API/internal/service"
	"TODO_API/pkg/database"
	"TODO_API/pkg/jwt"
	"TODO_API/pkg/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "TODO_API/docs" // 导入Swagger文档
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// 配置路由
func setupRouter(r *gin.Engine, h *handler.Healther, a *handler.AuthHandler, u *handler.UserHandeler, t *handler.TodoHandler) {
	// 添加Swagger文档路由（仅在开发环境）
	if config.GlobalConfig.App.Environment == "development" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	//健康，就绪检查路由
	r.GET("/health", h.HealthChecker)
	r.GET("/ready", h.ReadyChecker)

	//主页路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "欢迎使用 " + config.GlobalConfig.App.Name,
			"version":     config.GlobalConfig.App.Version,
			"environment": config.GlobalConfig.App.Environment,
			"endpoints": []string{
				"GET    /health - 健康检查",
				"POST   /api/auth/register - 用户注册",
				"POST   /api/auth/login - 用户登录",
				"POST   /api/auth/refresh - 刷新令牌",
				"GET    /api/users/me - 获取当前用户(需认证)",
				"GET    /api/todos - 获取待办事项列表(需认证)",
				"POST   /api/todos - 创建待办事项(需认证)",
				"GET    /api/todos/:id - 获取待办事项详情(需认证)",
				"PUT    /api/todos/:id - 更新待办事项(需认证)",
				"DELETE /api/todos/:id - 删除待办事项(需认证)",
				"PUT    /api/todos/:id/status - 更新状态(需认证)",
				"PUT    /api/todos/batch/status - 批量更新状态(需认证)",
			},
		})
	})

	//api路由组
	api := r.Group("/api")
	{
		//认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", a.Register)
			auth.POST("/login", a.Login)
			auth.POST("/refresh", a.RefreshToken)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleWare())
		{
			//用户相关路由
			user := protected.Group("/users")
			{
				user.GET("/me", u.GetProfile)
				user.PUT("/me", u.UpdateProfile)
				user.PUT("/me/password", u.ChangePassword)
			}

			// 待办事项路由
			todos := protected.Group("/todos")
			{
				// 基本CRUD操作
				todos.GET("", t.GetTodos)          // 获取待办事项列表
				todos.POST("", t.CreateTodo)       // 创建待办事项
				todos.GET("/:id", t.GetTodoByID)   // 获取单个待办事项
				todos.PUT("/:id", t.UpdateTodo)    // 更新待办事项
				todos.DELETE("/:id", t.DeleteTodo) // 删除待办事项

				// 状态操作
				todos.PUT("/:id/status", t.UpdateTodoStatus)    // 更新状态
				todos.PUT("/batch/status", t.BatchUpdateStatus) // 批量更新状态
			}
		}

	}

	// 调试路由：显示所有已注册的路由
	r.GET("/debug/routes", func(c *gin.Context) {
		var routes []map[string]string
		for _, route := range r.Routes() {
			routes = append(routes, map[string]string{
				"method": route.Method,
				"path":   route.Path,
			})
		}
		c.JSON(200, gin.H{
			"total":  len(routes),
			"routes": routes,
		})
	})
}

// 配置基础中间件
func setupBasicMiddleWare(g *gin.Engine) {
	//添加请求日志中间件
	g.Use(func(c *gin.Context) {
		logger.Info("Http请求",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)
		c.Next()
	})
	//恢复中间件
	g.Use(gin.Recovery())
}

// 启动服务器
func startSever(g *gin.Engine) {
	port := config.GlobalConfig.Server.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: g,
	}

	go func() {
		logger.Info("服务器启动信息",
			zap.String("address", "http://localhost:"+port),
			zap.String("mode", config.GlobalConfig.Server.Mode),
			zap.String("environment", config.GlobalConfig.App.Environment),
		)

		//log.Printf("健康检查: http://localhost:%s/health", port)
		//log.Printf("配置信息: http://localhost:%s/config", port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("服务器启动失败", zap.Error(err))
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭服务器")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	}
	logger.Info("服务器关闭成功")

}

// @title Go Todo API
// @version 1.0
// @description 基于 Go + Gin 的待办事项管理系统 API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 输入 "Bearer {token}" 进行认证
func main() {
	//初始化配置
	configPath := os.Getenv("CONFIG_PATH")
	config.InitConfig(configPath)

	//初始化JWT
	jwt.InitJWT()

	//初始化日志
	logger.InitLogger(config.GlobalConfig.Log.Level, config.GlobalConfig.Log.Filename)
	logger.Info("应用程序启动",
		zap.String("app", config.GlobalConfig.App.Name),
		zap.String("version", config.GlobalConfig.App.Version),
		zap.String("environment", config.GlobalConfig.App.Environment),
	)

	//连接数据库
	dbConfig := config.GlobalConfig.Database
	if err := database.InitMysql(dbConfig.Host, dbConfig.Port,
		dbConfig.User, dbConfig.Password, dbConfig.DBName); err != nil {
		logger.Error("数据库连接失败", zap.Error(err))
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer database.CloseMysql()

	//设置gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	//创建gin实例
	r := gin.Default()

	//添加基础中间件
	setupBasicMiddleWare(r)

	//初始化依赖注入
	userRepo := repository.NewUserRepository(database.GetDB())
	todoRepo := repository.NewTodoRepository(database.GetDB())

	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	todoService := service.NewTodoService(todoRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandeler(userService)
	todoHandler := handler.NewTodoHandler(todoService)
	healthHandler := handler.NewHealther()
	//设置路由
	setupRouter(r, healthHandler, authHandler, userHandler, todoHandler)

	//启动服务器
	startSever(r)
}
