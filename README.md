<h1 id="DvdnZ">GIN_</h1>
一个基于 Go + Gin 的后端服务，提供标准的分层架构与常用基础设施封装。

<h2 id="ulMZG">📁 项目结构</h2>
```plain
project/
├── api/                    # API 接口定义（业务服务接口）
├── cmd/
│   └── server/             # 服务启动入口
├── configs/                # 配置文件（YAML）
├── internal/
│   ├── handler/            # HTTP 请求处理器（Controller）
│   ├── model/              # 数据模型 & API 响应结构
│   ├── service/            # 业务逻辑实现
│   └── ...                 # （可扩展 repository/middleware 等）
├── pkg/
│   ├── config/             # 配置加载
│   ├── logger/             # 日志组件（Zap）
│   └── utils/              # 工具类（如 Snowflake ID）
└── go.mod
```

<h2 id="EdLn9">⚙️ 快速启动</h2>
<h3 id="VEysM">前提条件</h3>
+ Go 1.19+
+ Redis（用于部分功能）
+ MinIO 或兼容 S3 的对象存储（如需文件上传）

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

关键配置项说明（configs/app-dev.yaml）：

```yaml
# 服务配置
server:
  name: service_name
  port: 8080
  version: "1.0"

# 数据库配置
db:
  mysql:
    url: "root:aaaa@tcp(localhost:3306)/test?charset=utf8&parseTime=true&loc=Local"
    maxIdleConnection: 100
    maxOpenConnection: 130

  redis:
    addr: "127.0.0.1:6379"
    db: 0
    network: tcp
    username: ""
    password: "aaaa"

  tdengine:
    url: "root:aaaa@http(127.0.0.1:6041)/test"

# 消息队列配置
rabbitmq:
  url: "amqp://root:aaaa@127.0.0.1:5672/"

# 日志配置
log:
  path: "log/"
  level: debug
  maxSize: 500      # 单个日志文件最大大小（MB）
  maxBackups: 3     # 保留旧日志文件数量
  maxAge: 20        # 保留日志天数

# MinIO 对象存储配置
minio:
  enabled: true               
  bucketName: "test"
  region: "us-east-1"
  endpoint: "https://minio.org.cn/"
  accessKey: "XXXXXXXXXXXXXXXXXX"
  secretKey: "XXXXXXXXXXXXXXXXXXXXXX"

# 公共配置
public:
  ip: "http://127.0.0.1"    

# 调试配置
debug:
  enablePProf: false

# 定义启用的组件列表
components:
  - mysql
  # - redis
  # - minio
```

<h3 id="gZ4aK">4. 启动服务</h3>
```bash
go run cmd/server/main.go
```

服务默认运行在 [http://localhost:8080](http://localhost:8080)

---

<h2 id="OJwvX">🧪 示例接口</h2>
<h3 id="OIffD">获取数据</h3>
POST /api/data  
Content-Type: application/json

{  
    "param1": "value1",  
    "param2": "value2"  
}

成功响应：  
{  
    "code": 200,  
    "message": "success",  
    "data": {  
        // 业务数据  
    }  
}

错误响应：  
{  
    "code": 400,  
    "message": "参数解析出错"  
}

---

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

<h2 id="Ohh7x">🧪 测试</h2>
运行所有测试：

```bash
go test ./...
```

> 💡 建议为 service 层编写单元测试（通过 mock 接口）
>

---

<h2 id="eKAHz">🐳 Docker 部署（可选）</h2>
1. 构建镜像：

```bash
docker build -t project-backend .
```

2. 运行容器：

```bash
docker run -p 8080:8080 
  -v $(pwd)/configs:/app/configs 
    project-backend
```

---

