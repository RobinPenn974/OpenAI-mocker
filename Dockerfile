FROM golang:1.22.8 AS builder

WORKDIR /app

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o openai-mocker .

# 使用精简的alpine镜像来运行应用
FROM alpine:latest

WORKDIR /app

# 安装CA证书
RUN apk --no-cache add ca-certificates

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/openai-mocker .

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./openai-mocker"] 