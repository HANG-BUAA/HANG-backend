# 第一阶段：构建Go二进制文件
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将go.mod和go.sum复制到容器中
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 将代码复制到容器中
COPY ./src ./src

# 进入src目录并编译Go项目
WORKDIR /app/src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 第二阶段：创建一个更小的运行镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 加载配置文件
COPY ./src/config/settings_dev.yaml /root/config/settings.yaml
# 从构建阶段复制二进制文件
COPY --from=builder /app/src/main .

# 暴露服务端口
EXPOSE 8000

# 运行二进制文件
CMD ["./main"]
