# JWT 认证流程示例

## 介绍

### 使用已实现的部分

1. **JWT 工具包** (`pkg/jwt/jwt.go`)
   - Token 生成
   - Token 解析
   - Token 验证

2. **JWT 中间件** (`internal/middleware/auth.go`)
   - 请求拦截
   - Token 验证
   - 用户信息注入

3. **配置支持** (`configs/app-dev.yaml`)
   - Secret 密钥
   - Token 有效期
   - Refresh Token 有效期

### 关键部分

根据 JWT 最佳实践，你的项目缺少：

1. **认证端点**
   - 登录接口
   - 注册接口
   - Token 刷新接口
   - 登出接口

2. **Refresh Token 管理**
   - 数据库存储
   - Token 撤销机制
   - Token 轮换策略

3. **用户管理**
   - 用户表结构
   - 密码加密
   - 用户信息获取

### 目录结构

```
examples/jwt_auth/
├── model/              # 数据模型
│   ├── user.go
│   ├── auth_request.go
│   └── auth_response.go
├── repository/         # 数据访问层
│   ├── user_repository.go
│   └── token_repository.go
├── service/            # 业务逻辑层
│   └── auth_service.go
├── handler/            # HTTP 处理器
│   └── auth_handler.go
├── router/             # 路由注册
│   └── router.go
├── schema.sql          # 数据库表结构
└── README.md           # 使用文档
```

### 安全特性

1. **双 Token 机制**
   - Access Token：2 小时有效期
   - Refresh Token：7 天有效期

2. **Refresh Token 存储**
   - 数据库持久化
   - 可撤销设计
   - 支持多设备登录

3. **密码安全**
   - bcrypt 加密（cost=10）
   - 永不明文存储

4. **Token 轮换**
   - 刷新时撤销旧 token
   - 防止 token 泄露

5. **强制重登录场景**
   - 修改密码后
   - 登出所有设备

### 完整认证流程

```
1. 用户注册/登录
   ↓
2. 生成 Access Token (2h) + Refresh Token (7d)
   ↓
3. Refresh Token 存入数据库
   ↓
4. 返回给客户端
   ↓
5. 客户端使用 Access Token 访问 API
   ↓
6. Access Token 过期
   ↓
7. 使用 Refresh Token 刷新
   ↓
8. 撤销旧 Refresh Token，生成新的
   ↓
9. 返回新 Token Pair
```

### API 端点

#### 公开端点
- `POST /api/auth/register` - 注册
- `POST /api/auth/login` - 登录
- `POST /api/auth/refresh` - 刷新 token
- `POST /api/auth/logout` - 登出

#### 需要认证的端点
- `GET /api/auth/userinfo` - 获取用户信息
- `POST /api/auth/change-password` - 修改密码
- `POST /api/auth/logout-all` - 登出所有设备

## 使用示例

### 登录
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}'
```

### 访问受保护资源
```bash
curl -X GET http://localhost:8080/api/auth/userinfo \
  -H "Authorization: Bearer {access_token}"
```

### 刷新 Token
```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "your_refresh_token"}'
```

## 集成到主项目

在 `cmd/api/main.go` 中：

```go
import (
    authHandler "project/examples/jwt_auth/handler"
    authRepo "project/examples/jwt_auth/repository"
    authService "project/examples/jwt_auth/service"
    authRouter "project/examples/jwt_auth/router"
)

// 初始化
userRepo := authRepo.NewUserRepository(db)
tokenRepo := authRepo.NewTokenRepository(db)
authSvc := authService.NewAuthService(userRepo, tokenRepo, jwtConfig)
authHdl := authHandler.NewAuthHandler(authSvc)

// 注册路由
authRouter.RegisterAuthRoutes(r, authHdl)
```

## ⚠️ 生产环境注意事项

1. **使用强随机密钥**（至少 256 位）
2. **启用 HTTPS**
3. **Refresh Token 使用 httpOnly Cookie**
4. **定期清理过期 Token**
5. **记录审计日志**
6. **实现限流保护**