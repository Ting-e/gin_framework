#-------------- 构建阶段使用"builder"-------------------
    FROM golang:alpine AS builder

    # 设置容器时区（提前到构建阶段）
    RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
        && echo 'Asia/Shanghai' > /etc/timezone
    
    WORKDIR /builder
    # 优先复制依赖文件利用缓存，增加依赖安装管理
    COPY go.mod go.sum ./
    RUN go env -w GOPROXY=https://goproxy.cn,direct \
        && go mod download
    
    # 复制全部项目代码
    COPY . .
    #创建项目log目录
    RUN mkdir -p /builder/log
    
    # 构建应用（合并命令减少层数）
    RUN CGO_ENABLED=0 GOOS=linux go build -o main web/handler/main.go
    
    # -----------------生产阶段，使用更严格的指定alpine版本-------------
    FROM alpine:3.19 AS production
    
    # 更换镜像源
    RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories
    
    # 安装依赖
    RUN apk add --no-cache tzdata
    
    # 设置容器时区
    ENV TZ="Asia/Shanghai"
    RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime
    
    WORKDIR /app
    
    # 从构建阶段复制内容（优化路径结构）
    # 应用二进制
    COPY --from=builder /builder/main .  
    # 配置文件
    COPY --from=builder /builder/conf ./conf  
    # 日志目录
    COPY --from=builder /builder/log ./log   
    
    # 使用exec模式启动
    ENTRYPOINT ["./main"]
    CMD ["-conf=./conf/app-prod.yaml","-log=./log/"]
    # 