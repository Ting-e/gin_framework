# -------------- 构建阶段 -------------------
FROM golang:alpine AS builder

# 设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone

WORKDIR /builder

# 复制依赖并下载
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download

# 复制代码并构建
COPY . .

#创建项目log目录
RUN mkdir -p /builder/log

# 构建应用（合并命令减少层数
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# ----------------- 生产阶段 -----------------
FROM alpine:3.19 AS production

# 更换镜像源 + 安装依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories \
    && apk add --no-cache tzdata \
    && addgroup -g 10001 -S appgroup \
    && adduser -u 10001 -S appuser -G appgroup

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime

WORKDIR /app

# 从构建阶段复制二进制和配置、日志
COPY --from=builder /builder/main .
COPY --from=builder /builder/configs ./configs
COPY --from=builder /builder/log ./log

# 使用非 root 用户
USER appuser

ENTRYPOINT ["./main"]
CMD ["-conf=./configs/app-prod.yaml", "-log=./log/"]