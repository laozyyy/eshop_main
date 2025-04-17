package database_test

import (
	"context"
	"eshop_main/database"
	"testing"
)

func TestInitES(t *testing.T) {
	// 清理旧的ESClient引用
	database.ESClient = nil

	err := database.InitES()
	if err != nil {
		t.Fatalf("InitES失败: %v", err)
	}

	// 验证客户端是否创建成功
	if database.ESClient == nil {
		t.Fatal("ESClient未初始化")
	}

	// 测试Ping连接
	_, code, err := database.ESClient.Ping("http://117.72.72.114:19200").Do(context.Background())
	if err != nil || code != 200 {
		t.Fatalf("ES节点不可达，状态码:%d 错误:%v", code, err)
	}

	// 验证索引是否存在
	exists, err := database.ESClient.IndexExists("goods_index").Do(context.Background())
	if err != nil || !exists {
		t.Fatalf("索引检查失败 exists:%v 错误:%v", exists, err)
	}
}
