# 构建镜像
FROM golang:1.17.12-alpine AS builder

LABEL stage=gobuilder

# 设置镜像所需的必要环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct


# 移动到工作目录
WORKDIR /build

# 下载依赖
ADD go.mod .
ADD go.sum .
RUN go mod download

# 将源码复制到容器中
COPY . .

# 构建应用
RUN go build -ldflags="-s -w" -o /app/social

# 构建一个微镜像
FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

COPY --from=builder /app/social /bin/social

# 导出必须的端口
EXPOSE 1234

# 运行命令
CMD ["social"]