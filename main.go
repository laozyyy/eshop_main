package main

import (
	"eshop_main/database"
	"eshop_main/handler"
	"eshop_main/kitex_gen/eshop/home/goodsservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
)

func main() {
	database.Init()
	svr := goodsservice.NewServer(
		new(handler.GoodsServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "UserService"}),
	)
	err := svr.Run()
	if err != nil {
		panic(err)
	}
}
