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

## 🚀 快速开始

### 1. 前置条件

确保您的开发环境已安装：

- **Go** (1.20 或更高版本)
- **MySQL** (8.0 或更高版本)
- **Git**

### 2. 获取项目

```bash
git clone <您的仓库地址>
cd go-todo-api
```

### 3. 配置项目

1. **复制并修改配置文件**：

   ```bash
   cp config/config.yaml.example config/config.yaml
   # 使用编辑器打开 config/config.yaml，修改数据库连接等配置
   ```

2. **初始化数据库**：

   ```bash
   # 登录 MySQL，创建数据库
   # mysql -u root -p
   # CREATE DATABASE todo_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   
   # 执行初始化脚本，创建数据表
   mysql -u root -p todo_db < scripts/init_db.sql
   ```

### 4. 安装依赖 & 运行

```bash
# 下载项目依赖
go mod download

# 运行开发服务器 (默认端口 8080)
go run cmd/main.go
```

服务启动后，您可以在浏览器访问：

- **API 文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health

## 📖 API 使用指南

所有需要认证的 API 请求，都必须在 `Header`中携带：

```markdown
Authorization: Bearer <您的 access_token>
```

### 核心接口示例

#### 用户注册

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "confirm_password": "password123"
}
```

#### 用户登录

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "password123"
}
```

登录成功将返回 `access_token`和 `refresh_token`。

#### 创建待办事项

```http
POST /api/todos
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "学习 Go 语言",
  "description": "完成 Gin 框架的学习",
  "priority": 3,
  "due_date": "2024-12-31T23:59:59Z"
}
```

#### 获取待办事项列表 (支持分页和筛选)

```http
GET /api/todos?page=1&page_size=10&status=0&priority=3
Authorization: Bearer <access_token>
```

**更详细的接口定义和参数说明，请直接查看运行时的 Swagger 文档。**

## 🗄️ 数据库模型

项目核心包含以下数据表，其关系如下图所示：

- `users`: 用户表，存储账户信息。
- `todos`: 待办事项表，与用户关联。
- `tags`: 标签表，与用户关联。
- `todo_tags`: 待办事项与标签的关联表（多对多关系）。

完整的 `CREATE TABLE`SQL 语句请查看 `scripts/init_db.sql`文件。

## 🔧 配置说明

核心配置位于 `config/config.yaml`，主要部分如下：

```yaml
server:
  port: "8080"          # 服务端口
  mode: "debug"         # 运行模式: debug 或 release

database:
  host: "localhost"
  port: "3306"
  user: "root"
  password: "your_password"
  dbname: "todo_db"
  charset: "utf8mb4"

jwt:
  secret: "your-secret-key-change-in-production" # 请务必在生产环境中更改！
  access_expire: 3600    # 访问令牌有效期(秒)
  refresh_expire: 604800 # 刷新令牌有效期(秒)

log:
  level: "debug"         # 日志级别: debug, info, warn, error
  filename: "logs/app.log"
```