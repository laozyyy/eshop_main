FROM golang:latest

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖并缓存
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o eshop_main

EXPOSE 8888
ENTRYPOINT ["./eshop_main"]