package database

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// 新增ES文档结构体（只包含mapping中定义的字段）
type ESGoodsDocument struct {
	Sku       string  `json:"sku"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Spec      string  `json:"spec"`
	TagID     string  `json:"tag_id"`
	IsDeleted int32   `json:"is_deleted"`
}

func TestSyncAllGoodsToES(t *testing.T) {
	// 初始化数据库
	Init()

	// 获取所有未删除的商品
	var goods []GoodsSku
	if err := DB.Where("is_deleted = 0").Find(&goods).Error; err != nil {
		t.Fatalf("数据库查询失败: %v", err)
	}

	success := 0
	for i, g := range goods {
		// 使用ES客户端同步数据
		// 转换数据库结构体为ES专用结构体
		doc := ESGoodsDocument{
			Sku:       g.Sku,
			Name:      g.Name,
			Price:     g.Price,
			Spec:      g.Spec,
			TagID:     g.TagID,
			IsDeleted: g.IsDeleted,
		}
		fmt.Printf("第%d\n个", i)
		_, err := ESClient.Index().
			Index("goods_index").
			Id(g.Sku).
			BodyJson(doc). // 使用过滤后的结构体
			Refresh("wait_for").
			Do(context.Background())

		if err != nil {
			t.Errorf("同步失败 SKU:%s 错误:%v", g.Sku, err)
			continue
		}
		success++
	}

	t.Logf("同步完成 总数:%d 成功:%d 失败:%d",
		len(goods),
		success,
		len(goods)-success)
}

func TestESSearchGoods(t *testing.T) {
	// 初始化服务
	Init()

	// 测试数据准备
	testKeyword := "阿" // 可修改为实际存在的关键词
	pageSize := int32(100)
	pageNum := int32(1)

	// 执行搜索
	skus, isEnd, err := SearchGoodsByName(testKeyword, pageSize, pageNum)
	if err != nil {
		t.Fatalf("ES查询失败: %v", err)
	}

	// 验证基础结果
	if len(skus) == 0 {
		t.Fatal("未查询到任何商品")
	}

	// 验证分页参数
	if len(skus) > int(pageSize) {
		t.Errorf("返回结果数量超过分页限制 预期:%d 实际:%d", pageSize, len(skus))
	}

	// 验证字段匹配
	firstSku := skus[0]
	if firstSku.Name == "" {
		t.Error("商品名称字段为空")
	}
	if firstSku.Price <= 0 {
		t.Error("商品价格字段异常")
	}
	for _, sku := range skus {
		fmt.Println(sku.Name)
	}

	// 验证搜索关键词匹配
	if !strings.Contains(strings.ToLower(firstSku.Name), strings.ToLower(testKeyword)) {
		t.Errorf("商品名称不包含搜索关键词 名称:%s 关键词:%s", firstSku.Name, testKeyword)
	}

	t.Logf("搜索完成 关键词:%s 结果数量:%d 是否结束:%v",
		testKeyword,
		len(skus),
		isEnd)
}
