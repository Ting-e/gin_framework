<h1 id="DvdnZ">GIN_framework</h1>
一个基于 Go + Gin 的后端服务，提供标准的分层架构与常用基础设施封装。

<h2 id="EdLn9">⚙️ 快速启动</h2>
<h3 id="VEysM">前提条件</h3>

+ Go 1.19+

<h3 id="mqbE7">1. 克隆项目</h3>

```bash
git clone https://github.com/Ting-e/gin_framework.git
cd project
```


<h3 id="UBpsU">2. 安装依赖</h3>

```bash
go mod tidy
```

<h3 id="SZjam">3. 配置环境</h3>

复制配置模板并按需修改：

```bash
cp configs/app-dev.yaml.example configs/app-dev.yaml
```

<h3 id="gZ4aK">4. 启动服务</h3>

```bash
go run cmd/server/main.go
```

服务默认运行在 [http://localhost:8080](http://localhost:8080)

<h2 id="DwJ9l">🛠 核心功能</h2>

+ ✅ 标准化分层：Handler → Service → Repository
  
+ ✅ 统一响应格式：所有接口返回 {code, message, data}
  
+ ✅ 结构化日志：使用 Zap 记录请求/响应及错误
  
+ ✅ ID 生成：集成 Snowflake 分布式 ID
  
+ ✅ 文件上传：
  
    - 后端直传 MinIO  
    - 预签名 URL 上传（前端直传）
      
+ ✅ HTTP 客户端封装：简化第三方服务调用  

---

<h2 id="eKAHz">🐳 Docker 部署（可选）</h2>

构建镜像：


```bash
docker build -t project-backend .
```

