# TODO API 🚀

一个基于 **Go + Gin + GORM** 构建的、功能完整的待办事项（Todo）管理后端 API。

## ✨ 特性

- **🔐 完整的用户认证系统**: 使用 JWT 实现安全的用户注册、登录和令牌刷新机制。
- **✅ 核心业务功能**: 提供待办事项（Todo）的增删改查（CRUD）、状态管理。
- **📖 自动化 API 文档**: 集成 Swagger，自动生成并在线提供交互式 API 文档。
- **🔧 多环境配置**: 通过 Viper 支持开发、生产等多套环境配置。
- **📝 结构化日志**: 使用高性能的 Zap 日志库进行结构化日志记录。
- **🧪 清晰的架构**: 采用分层架构（Handler -> Service -> Repository），代码结构清晰，易于维护和测试。
- **⚙️ 统一响应与错误处理**: 标准化的 API 响应格式和完善的错误处理机制。
- **🔄 缓存机制**: 集成 Redis 缓存，优化用户数据查询性能。
- **📡 消息队列**: 集成 RabbitMQ 消息队列，实现事件驱动架构。

## 🛠️ 技术栈

| 组件                 | 说明               | 版本  |
| -------------------- | ------------------ | ----- |
| **Go**               | 后端编程语言       | 1.25.0 |
| **Gin**              | 高性能 Web 框架    | 1.11.0 |
| **GORM**             | 功能强大的 ORM 库  | 1.31.1 |
| **MySQL**            | 关系型数据库       | 8.0+  |
| **Redis**            | 缓存数据库         | 最新  |
| **RabbitMQ**         | 消息队列           | 最新  |
| **JWT (golang-jwt)** | 用户身份认证       | 4.5.2 |
| **Viper**            | 配置管理工具       | 1.21.0 |
| **Zap**              | 高性能日志库       | 1.27.1 |
| **Swagger**          | API 文档生成与展示 | 1.16.6 |

## 📁 项目结构

```markdown
TODO_API/
├── cmd/main.go                 # 应用程序入口
├── config/                     # 配置文件目录
│   ├── config.go              # 配置结构体定义
│   ├── config.yaml            # 主配置文件
│   ├── config.dev.yaml        # 开发环境配置
│   └── config.prod.yaml       # 生产环境配置
├── internal/                   # 内部应用代码
│   ├── app/                   # 应用层
│   │   ├── dto/               # 数据传输对象 (请求/响应结构)
│   │   │   ├── request/       # 请求结构体
│   │   │   └── response/      # 响应结构体
│   │   ├── handler/           # HTTP 请求处理器 (类似 Controller)
│   │   └── middleware/        # Gin 中间件 (认证、日志、跨域等)
│   ├── domain/                # 领域层
│   │   ├── event/             # 事件定义
│   │   └── model/             # 数据库模型 (GORM 结构体)
│   ├── infrastructure/        # 基础设施层
│   │   └── persistence/       # 持久化相关
│   ├── repository/            # 数据访问层 (数据库操作)
│   └── service/               # 业务逻辑层
├── pkg/                       # 公共库包
│   ├── cache/                 # 缓存相关
│   ├── database/              # 数据库连接初始化
│   ├── encryption/            # 加密相关
│   ├── errors/                # 错误处理
│   ├── jwt/                   # JWT 令牌工具
│   ├── lock/                  # 分布式锁
│   ├── logger/                # 日志初始化
│   ├── mq/                    # 消息队列
│   ├── ratelimit/             # 速率限制
│   ├── response/              # 统一 API 响应格式
│   ├── storage/               # 文件存储
│   ├── token/                 # 令牌管理
│   └── validator/             # 数据验证
├── scripts/                   # 脚本目录
│   └── init_db.sql           # 数据库初始化 SQL 脚本
├── docs/                      # 项目文档
├── go.mod                     # Go 模块定义文件
├── go.sum                     # Go 依赖校验文件
├── Dockerfile                 # Docker 构建文件
├── docker-compose.yml         # Docker Compose 配置文件
└── README.md                  # 项目说明文档 (本文件)
```

## 🚀 快速开始

### 环境要求

- Go 1.25.0+
- MySQL 8.0+
- Redis (可选，用于缓存和令牌管理)
- RabbitMQ (可选，用于消息队列)

### 安装与运行

1. **克隆项目**

```bash
git clone https://github.com/Bulonte/TODO_API.git
cd TODO_API
```

2. **配置环境**

复制配置文件并修改为你的环境配置：

```bash
cp config/config.yaml config/config.local.yaml
# 编辑 config.local.yaml 文件，设置数据库连接等信息
```

3. **初始化数据库**

执行 SQL 脚本初始化数据库：

```bash
mysql -u root -p < scripts/init_db.sql
```

4. **安装依赖**

```bash
go mod tidy
```

5. **运行项目**

```bash
# 开发环境
go run cmd/main.go

# 生产环境
CONFIG_PATH=config/config.prod.yaml go run cmd/main.go
```

### Docker 部署

使用 Docker Compose 快速部署：

```bash
docker-compose up -d
```

## 📖 API 文档

项目集成了 Swagger 文档，在开发环境下可以通过以下地址访问：

```
http://localhost:8080/swagger/index.html
```

## 🔗 API 端点

### 认证相关

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/refresh` - 刷新令牌

### 用户相关

- `GET /api/users/me` - 获取当前用户信息
- `PUT /api/users/me` - 更新用户信息
- `PUT /api/users/me/password` - 修改密码

### 待办事项相关

- `GET /api/todos` - 获取待办事项列表
- `POST /api/todos` - 创建待办事项
- `GET /api/todos/:id` - 获取单个待办事项
- `PUT /api/todos/:id` - 更新待办事项
- `DELETE /api/todos/:id` - 删除待办事项
- `PUT /api/todos/:id/status` - 更新待办事项状态
- `PUT /api/todos/batch-status` - 批量更新待办事项状态

## 📝 项目配置

项目使用 Viper 进行配置管理，支持多个环境的配置文件：

- `config.yaml` - 主配置文件
- `config.dev.yaml` - 开发环境配置
- `config.prod.yaml` - 生产环境配置

主要配置项包括：

- 应用配置（名称、版本、环境）
- 服务器配置（端口、模式）
- 数据库配置（连接信息）
- Redis 配置（连接信息）
- RabbitMQ 配置（连接信息）
- JWT 配置（密钥、过期时间）
- 日志配置（级别、文件路径）
- 存储配置（基础路径、基础 URL）

