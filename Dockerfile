FROM golang:latest

WORKDIR /app/demo
COPY . .

RUN go build eshop_main

EXPOSE 8888
ENTRYPOINT ["./eshop_main"]