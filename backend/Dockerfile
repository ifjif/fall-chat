FROM docker.m.daocloud.io/library/golang:1.25-alpine

WORKDIR /app

# 设置 GOPROXY 环境变量，使用国内代理加速下载
ENV GOPROXY=https://goproxy.cn,direct
#安装热重载工具 air
RUN go install github.com/air-verse/air@v1.67.1

COPY go.mod go.sum .

RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]
