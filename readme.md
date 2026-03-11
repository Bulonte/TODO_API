# Go-Todo-API 🚀

一个基于 **Go + Gin + GORM** 构建的、功能完整的待办事项（Todo）管理后端 API。

## ✨ 特性

- **🔐 完整的用户认证系统**: 使用 JWT 实现安全的用户注册、登录和令牌刷新机制。
- **✅ 核心业务功能**: 提供待办事项（Todo）的增删改查（CRUD）、状态与优先级管理。
- **🏷️ 灵活的标签系统**: 支持为待办事项创建、管理标签，实现分类与筛选。
- **📖 自动化 API 文档**: 集成 Swagger，自动生成并在线提供交互式 API 文档。
- **🔧 多环境配置**: 通过 Viper 支持开发、生产等多套环境配置。
- **📝 结构化日志**: 使用高性能的 Zap 日志库进行结构化日志记录。
- **🧪 清晰的架构**: 采用分层架构（Handler -> Service -> Repository），代码结构清晰，易于维护和测试。
- **⚙️ 统一响应与错误处理**: 标准化的 API 响应格式和完善的错误处理机制。

## 🛠️ 技术栈

| 组件                 | 说明               | 版本  |
| -------------------- | ------------------ | ----- |
| **Go**               | 后端编程语言       | 1.20+ |
| **Gin**              | 高性能 Web 框架    | 最新  |
| **GORM**             | 功能强大的 ORM 库  | 最新  |
| **MySQL**            | 关系型数据库       | 8.0+  |
| **JWT (golang-jwt)** | 用户身份认证       | 最新  |
| **Viper**            | 配置管理工具       | 最新  |
| **Zap**              | 高性能日志库       | 最新  |
| **Swagger**          | API 文档生成与展示 | 最新  |

## 📁 项目结构

```markdown
go-todo-api/
├── cmd/main.go                 # 应用程序入口
├── config/                     # 配置文件目录
│   ├── config.go              # 配置结构体定义
│   ├── config.yaml            # 主配置文件
│   ├── config.dev.yaml        # 开发环境配置
│   └── config.prod.yaml       # 生产环境配置
├── internal/                   # 内部应用代码
│   ├── app/                   # 应用层
│   │   ├── dto/               # 数据传输对象 (请求/响应结构)
│   │   ├── handler/           # HTTP 请求处理器 (类似 Controller)
│   │   └── middleware/        # Gin 中间件 (认证、日志、跨域等)
│   ├── domain/                # 领域层
│   │   └── model/             # 数据库模型 (GORM 结构体)
│   ├── repository/            # 数据访问层 (数据库操作)
│   └── service/               # 业务逻辑层
├── pkg/                       # 公共库包
│   ├── database/              # 数据库连接初始化
│   ├── jwt/                   # JWT 令牌工具
│   ├── logger/                # 日志初始化
│   └── response/              # 统一 API 响应格式
├── scripts/                   # 脚本目录
│   └── init_db.sql           # 数据库初始化 SQL 脚本
├── api/                       # API 描述文件 (Swagger)
├── docs/                      # 项目文档
├── go.mod                     # Go 模块定义文件
└── README.md                  # 项目说明文档 (本文件)
```

