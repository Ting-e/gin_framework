# 快速开始指南

## 5 分钟上手

### 第一步：安装依赖

```bash
# 克隆项目
git clone https://github.com/yourusername/backend_framework.git
cd backend_framework

# 下载依赖
go mod download
```

### 第二步：配置数据库

1. 创建 MySQL 数据库：
```sql
CREATE DATABASE test_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 修改配置文件 `configs/app-dev.yaml`：
```yaml
mysql:
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "your_password"  # 改成你的密码
  database: "test_db"
```

### 第三步：运行项目

```bash
# 方式1：使用 go run
go run cmd/api/main.go -conf ./configs/app-dev.yaml

# 方式2：使用 Makefile
make run
```

### 第四步：测试接口

访问健康检查接口：
```bash
curl http://localhost:8080/health
```

应该看到：
```json
{
  "status": "ok",
  "version": "1.0.0",
  "name": "Backend Framework"
}
```

---

## 第一个 API

### 1. 创建用户表

```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 插入测试数据（密码是 123456 的 bcrypt hash）
INSERT INTO users (username, password_hash, role) VALUES 
('admin', '$2a$14$xxxxx', 'admin');
```

### 2. 测试登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'
```

响应：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 7200
  }
}
```

### 3. 使用 Token 访问受保护接口

```bash
TOKEN="your_access_token_here"

curl -X GET http://localhost:8080/api/auth/userinfo \
  -H "Authorization: Bearer $TOKEN"
```

---

## 常用命令

```bash
# 运行
make run

# 编译
make build

# 测试
make test

# 清理
make clean

# Docker 运行
make docker-build
make docker-run
```

---

## 下一步

- 📖 阅读完整 [README.md](README.md)
- 🔍 查看 [examples/](examples/) 目录的示例代码
- 🛠️ 学习如何[添加新功能](README.md#开发指南)
- 📡 参考 [API 文档](README.md#api-文档)

---

## 遇到问题？

1. 检查配置文件是否正确
2. 确认数据库连接正常
3. 查看日志文件 `log/` 目录
4. 提交 Issue：https://github.com/Ting-e/gin_framework/issues