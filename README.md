# Backend Framework

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)

一个现代化的 Go 语言 Web 后端框架，基于 Gin，提供开箱即用的企业级功能。

[特性](#特性) • [快速开始](#快速开始) • [文档](#项目结构) • [示例](#示例代码) • [贡献](#贡献指南)

</div>

---

## 📋 目录

- [特性](#特性)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [项目结构](#项目结构)
- [配置说明](#配置说明)
- [示例代码](#示例代码)
- [中间件](#中间件)
- [数据库](#数据库)
- [API 文档](#api-文档)
- [部署](#部署)
- [开发指南](#开发指南)
- [常见问题](#常见问题)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

---

## ✨ 特性

### 核心功能

- 🚀 **高性能**：基于 Gin 框架，性能优异
- 🔐 **JWT 认证**：完整的 JWT 认证和授权系统
- 🗄️ **多数据库支持**：MySQL、Redis、TDengine
- 📝 **日志系统**：基于 Zap 的结构化日志，自动日志轮转
- 🔧 **配置管理**：YAML 配置文件，支持环境变量
- 🛡️ **安全中间件**：CORS、限流、超时控制、Panic 恢复
- 📦 **组件化设计**：按需加载数据库、消息队列等组件
- 🐰 **消息队列**：RabbitMQ 集成
- 📊 **性能分析**：内置 PProf 性能分析工具
- 🎯 **标准分层**：Handler → Service → Repository 架构

### 工具包

- **Response**：统一的 HTTP 响应格式
- **JWT**：Token 生成、验证、刷新
- **Snowflake**：分布式 ID 生成器
- **Storage**：文件上传下载（支持本地、OSS）
- **HTTP Client**：封装的 HTTP 请求工具

---

## 🛠️ 技术栈

| 类别 | 技术 |
|------|------|
| **框架** | [Gin](https://github.com/gin-gonic/gin) |
| **ORM** | [GORM](https://gorm.io/) |
| **日志** | [Zap](https://github.com/uber-go/zap) |
| **配置** | [Viper](https://github.com/spf13/viper) / YAML |
| **JWT** | [golang-jwt](https://github.com/golang-jwt/jwt) |
| **数据库** | MySQL, Redis, TDengine |
| **消息队列** | RabbitMQ |

---

## 🚀 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+ / 8.0+
- Redis 5.0+

### 安装

```bash
# 克隆项目
git clone https://github.com/Ting-e/gin_framework.git
cd gin_framework

# 下载依赖
go mod download
```

### 配置

复制并修改配置文件：

```bash
cp configs/app-dev.yaml configs/app-local.yaml
# 编辑 app-local.yaml，修改数据库连接等配置
```

### 运行

```bash
# 开发环境运行
go run cmd/api/main.go -conf ./configs/app-dev.yaml

# 或使用 Makefile
make run

# 指定配置文件和日志目录
go run cmd/api/main.go -conf ./configs/app-prod.yaml -log ./logs/
```

### 验证

访问 http://localhost:8080/，应该看到：

```json
{
  "ServiceName": "Test-server",
  "Now": "2025-11-27 18:15:50",
  "Environment": "Development",
  "Version": "v1.0.22",
  "StatusMessage": "Test-server is running"
}
```

---

## 📁 项目结构（以实际项目结构为准）

```
backend_framework/
├── cmd/                        # 应用入口
│   └── api/
│       └── main.go            # 主程序入口
├── configs/                    # 配置文件
│   ├── app-dev.yaml           # 开发环境配置
│   └── app-prod.yaml          # 生产环境配置
├── internal/                   # 内部应用代码
│   ├── app/                   # 应用初始化
│   │   ├── app.go             # 应用核心逻辑
│   │   ├── .go             # 应用核心逻辑
│   │   ├── config_loader.go   # 配置加载
│   │   └── logo.go            # Logo 显示
│   ├── handler/               # HTTP 处理器（Controller）
│   ├── service/               # 业务逻辑层
│   ├── repository/            # 数据访问层
│   ├── model/                 # 数据模型（DTO/VO）
│   ├── middleware/            # 中间件
│   │   ├── logger.go          # 日志中间件
│   │   ├── cors.go            # CORS 中间件
│   │   ├── recovery.go        # Panic 恢复
│   │   ├── auth.go            # JWT 认证
│   │   ├── rate_limit.go      # 限流
│   │   └── timeout.go         # 超时控制
│   └── router/                # 路由注册
├── pkg/                        # 可复用的公共包
│   ├── config/                # 配置管理
│   ├── database/              # 数据库连接
│   │   ├── mysql.go           # MySQL（原生）
│   │   ├── gorm.go            # GORM
│   │   ├── redis.go           # Redis
│   │   └── tdengine.go        # TDengine
│   ├── logger/                # 日志工具
│   ├── response/              # 统一响应
│   ├── jwt/                   # JWT 工具
│   ├── queue/                 # 消息队列
│   │   └── rabbitmq.go
│   └── utils/                 # 工具函数
│       ├── idgen/             # ID 生成器
│       ├── snowflake/         # 雪花算法
│       ├── storage/           # 文件存储
│       └── httpclient/        # HTTP 客户端
├── examples/                   # 示例代码
│   ├── simple_crud/           # 原生 SQL CRUD 示例
│   └── gorm_crud/             # GORM CRUD 示例
├── web/                        # 前端资源
├── log/                        # 日志目录
├── Dockerfile                  # Docker 构建文件
├── go.mod                      # Go 模块定义
├── go.sum                      # 依赖锁定
├── Makefile                    # 编译脚本
└── README.md                   # 项目文档
```

### 架构分层说明

```
┌─────────────────────────────────────────┐
│            HTTP Request                  │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│          Middleware Layer                │
│  (CORS, Auth, Logger, RateLimit...)     │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Handler Layer (控制器)           │
│  - 参数验证                              │
│  - 调用 Service                          │
│  - 返回响应                              │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Service Layer (业务层)           │
│  - 业务逻辑处理                          │
│  - 事务控制                              │
│  - 调用 Repository                       │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│       Repository Layer (数据层)          │
│  - 数据库操作                            │
│  - 缓存操作                              │
│  - 外部服务调用                          │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│          Database / Cache                │
└─────────────────────────────────────────┘
```

---

## ⚙️ 配置说明

### 配置文件示例 (configs/app-dev.yaml)

```yaml
# 服务器配置
server:
  name: "Backend Framework"
  port: 8080
  version: "1.0.0"

# 日志配置
log:
  level: "debug"  # debug, info, warn, error

# JWT 配置
jwt:
  secret: "your-256-bit-secret-key-change-in-production"
  issuer: "backend-framework"
  expires_hours: 2           # 访问令牌有效期（小时）
  refresh_expires_days: 7    # 刷新令牌有效期（天）

# MySQL 配置
mysql:
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "password"
  database: "test_db"
  charset: "utf8mb4"
  max_idle_conns: 10
  max_open_conns: 100
  max_lifetime: 3600

# Redis 配置
redis:
  host: "127.0.0.1"
  port: 6379
  password: ""
  db: 0

# TDengine 配置
tdengine:
  host: "127.0.0.1"
  port: 6041
  username: "root"
  password: "taosdata"
  database: "test"

# RabbitMQ 配置
rabbitmq:
  host: "127.0.0.1"
  port: 5672
  username: "guest"
  password: "guest"
  vhost: "/"

# 组件列表（按需加载）
components:
  - mysql
  - redis
  # - gorm
  # - tdengine
  # - rabbitmq

# 调试配置
debug:
  enable_pprof: true  # 性能分析工具
```

### 环境变量

支持通过环境变量覆盖配置：

```bash
# 配置文件路径
export APP_CONFIG=/path/to/config.yaml

# 日志目录
export APP_LOG_DIR=/path/to/logs

# 运行
go run cmd/api/main.go
```

---

## 💡 示例代码

### 1. 原生 SQL CRUD 示例

位置：`examples/simple_crud/`

**特点**：
- 使用原生 SQL 语句
- 适合复杂查询和性能优化
- 完整的增删改查示例

**目录结构**：
```
examples/simple_crud/
├── handler/        # HTTP 处理器
├── service/        # 业务逻辑
├── repository/     # 数据访问
├── model/          # 数据模型
└── router/         # 路由注册
```

### 2. GORM CRUD 示例

位置：`examples/gorm_crud/`

**特点**：
- 使用 GORM ORM
- 快速开发
- 支持关联查询、事务

### 示例：用户登录

```go
// handler/auth.go
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "参数错误")
        return
    }

    // 调用 service
    token, err := h.authService.Login(req.Username, req.Password)
    if err != nil {
        response.BusinessError(c, errcode.UserPasswordError, err.Error())
        return
    }

    response.Success(c, gin.H{"token": token})
}

// service/auth.go
func (s *AuthService) Login(username, password string) (string, error) {
    // 从 repository 获取用户
    user, err := s.userRepo.GetByUsername(username)
    if err != nil {
        return "", errors.New("用户不存在")
    }

    // 验证密码
    if !CheckPassword(password, user.PasswordHash) {
        return "", errors.New("密码错误")
    }

    // 生成 JWT
    token, _ := jwt.GenerateToken(user.ID, user.Username, user.Role)
    return token, nil
}

// repository/user.go
func (r *UserRepository) GetByUsername(username string) (*User, error) {
    var user User
    err := r.db.QueryRow(
        "SELECT id, username, password_hash, role FROM users WHERE username = ?",
        username,
    ).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
    return &user, err
}
```

---

## 🛡️ 中间件

### 内置中间件

| 中间件 | 说明 | 位置 |
|--------|------|------|
| **Logger** | 请求日志记录 | `internal/middleware/logger.go` |
| **CORS** | 跨域处理 | `internal/middleware/cors.go` |
| **Recovery** | Panic 恢复 | `internal/middleware/recovery.go` |
| **JWTAuth** | JWT 认证 | `internal/middleware/auth.go` |
| **RateLimit** | 限流 | `internal/middleware/rate_limit.go` |
| **Timeout** | 请求超时 | `internal/middleware/timeout.go` |

### 使用示例

```go
// 全局中间件
r := gin.New()
middleware.RegisterDefaultMiddlewares(r, version)

// 路由组中间件
auth := r.Group("/api")
auth.Use(middleware.JWTAuth())
{
    auth.GET("/profile", handler.GetProfile)
}

// 特定路由中间件
r.GET("/admin/users", 
    middleware.JWTAuth(),
    middleware.RequireRole("admin"),
    handler.GetUsers,
)
```

---

## 🗄️ 数据库

### MySQL

#### 原生 SQL

```go
import "project/pkg/database"

// 获取连接
db := database.GetMysql().GetDB()

// 查询
rows, err := db.Query("SELECT * FROM users WHERE id = ?", userID)

// 插入
result, err := db.Exec("INSERT INTO users (username) VALUES (?)", username)
```

#### GORM

```go
import "project/pkg/database"

// 获取 GORM 实例
db := database.GetGormMysql().GetDB()

// 查询
var user User
db.First(&user, "username = ?", "admin")

// 创建
db.Create(&User{Username: "test"})

// 更新
db.Model(&user).Update("email", "test@example.com")

// 删除
db.Delete(&user)
```

### Redis

```go
import "project/pkg/database"

// 获取 Redis 客户端
rdb := database.GetRedis().GetClient()

// 设置值
rdb.Set(ctx, "key", "value", 5*time.Minute)

// 获取值
val, err := rdb.Get(ctx, "key").Result()

// 删除
rdb.Del(ctx, "key")
```

---

<!-- ## 📡 API 文档

### 认证接口

#### 登录

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
```

**响应**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 7200,
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin"
    }
  }
}
```

#### 获取用户信息

```http
GET /api/auth/userinfo
Authorization: Bearer {access_token}
```

**响应**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

### 统一响应格式

#### 成功响应

```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

#### 错误响应

```json
{
  "code": 400,
  "message": "参数错误"
}
```

#### 分页响应

```json
{
  "code": 200,
  "message": "success",
  "data": [ ... ],
  "total": 100,
  "page": 1,
  "size": 10
}
```

--- -->

## 🚢 部署

### Docker 部署

```bash
# 构建镜像
docker build -t backend-framework:latest .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/log:/app/log \
  -e APP_CONFIG=/app/configs/app-prod.yaml \
  --name backend-api \
  backend-framework:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./configs:/app/configs
      - ./log:/app/log
    environment:
      - APP_CONFIG=/app/configs/app-prod.yaml
    depends_on:
      - mysql
      - redis

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: test_db
    ports:
      - "3306:3306"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
```

### 编译部署

```bash
# 编译
go build

# 运行
./backend_framework -conf ./configs/app-prod.yaml -log ./log/
```

---

## 🔨 开发指南

### 添加新功能

#### 1. 定义模型 (internal/model/)

```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3"`
    Password string `json:"password" binding:"required,min=6"`
}
```

#### 2. 数据访问层 (internal/repository/)

```go
func (r *UserRepository) Create(user *User) error {
    return r.db.Create(user).Error
}
```

#### 3. 业务逻辑层 (internal/service/)

```go
func (s *UserService) CreateUser(req *CreateUserRequest) error {
    user := &User{
        Username: req.Username,
        Password: HashPassword(req.Password),
    }
    return s.userRepo.Create(user)
}
```

#### 4. 控制器 (internal/handler/)

```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, err.Error())
        return
    }
    
    if err := h.userService.CreateUser(&req); err != nil {
        response.ServerError(c, err.Error())
        return
    }
    
    response.Success(c, nil)
}
```

#### 5. 注册路由 (internal/router/)

```go
func RegisterRoutes(r *gin.Engine) {
    api := r.Group("/api")
    api.POST("/users", handler.CreateUser)
}
```

### 代码规范

- 使用 `gofmt` 格式化代码
- 遵循 Go 官方编码规范
- 变量命名使用驼峰式
- 接口定义在实现之前
- 错误处理不能省略

### 测试

```bash
# 运行所有测试
go test ./...

# 测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## ❓ 常见问题

### Q1: 如何切换数据库？

修改 `configs/app-dev.yaml` 中的 `components` 列表：

```yaml
components:
  - mysql    # 原生 SQL
  # - gorm   # GORM
  - redis
```

### Q2: JWT Token 过期后怎么办？

使用 Refresh Token 刷新：

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your_refresh_token"
}
```

### Q3: 如何添加自定义中间件？

在 `internal/middleware/` 创建文件，然后在 `middleware.go` 注册：

```go
func RegisterDefaultMiddlewares(engine *gin.Engine, version string) {
    engine.Use(Recovery())
    engine.Use(Logger())
    engine.Use(YourCustomMiddleware())  // 添加这里
}
```

### Q4: 日志文件太大怎么办？

日志自动按天分割，配置在 `pkg/logger/logger.go`：

```go
MaxSize:    500,  // MB
MaxBackups: 3,    // 保留文件数
MaxAge:     28,   // 天
```

---

## 🤝 贡献指南

欢迎贡献代码！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

### 提交规范

```
feat: 新功能
fix: 修复 Bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试相关
chore: 构建/工具链相关
```

---

## 📄 许可证

本项目采用 [MIT License](LICENSE.txt) 开源协议。

---

## 📞 联系方式

- 作者：Tinge
- QQ群：1067520714
- 项目地址：https://github.com/Ting-e/gin_framework

---

## 🙏 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Zap](https://github.com/uber-go/zap)
- [golang-jwt](https://github.com/golang-jwt/jwt)

---

<div align="center">

**如果这个项目对您有帮助，请给一个 ⭐️ Star！**

Made with ❤️ by Your Name

</div>