FROM golang:alpine AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o fabric-edgenode

FROM debian:bullseye-slim

# 从builder镜像中把配置文件拷贝到当前目录
COPY ./cfg /conf
COPY ./kafka_crypto /kafka_crypto
COPY ./fixtures /fixtures

# 从builder镜像中把二进制文件拷贝到当前目录
COPY --from=builder /build/fabric-edgenode /

# 启动服务
CMD ["/fabric-edgenode"]
