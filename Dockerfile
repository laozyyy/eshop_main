FROM golang:latest

WORKDIR /app/demo
COPY . .

RUN go build eshop_api

EXPOSE 8888
ENTRYPOINT ["./eshop_api"]