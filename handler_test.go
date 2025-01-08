package main

import (
	"context"
	"eshop_main/database"
	"eshop_main/kitex_gen/eshop/home"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func getTestDB() *gorm.DB {
	dsn := "root:123456@tcp(117.72.72.114:13306)/eshop_main?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func TestGetOneSku(t *testing.T) {
	impl := &GoodsServiceImpl{}

	// 测试存在的SKU
	resp, err := impl.GetOneSku(context.Background(), "261311")
	if err != nil {
		t.Errorf("GetOneSku error: %v", err)
		return
	}
	fmt.Printf("GetOneSku response: %+v\n", resp)

	// 测试不存在的SKU
	resp, err = impl.GetOneSku(context.Background(), "not_exist_sku")
	if err != nil {
		t.Errorf("GetOneSku error: %v", err)
		return
	}
	fmt.Printf("GetOneSku not exist response: %+v\n", resp)
}

func TestMGetSku(t *testing.T) {
	impl := &GoodsServiceImpl{}

	// 测试分页查询
	req := &home.MGetSkuRequest{
		PageSize: 10,
		PageNum:  2,
		TagId:    "N_KjVQ",
	}

	resp, err := impl.MGetSku(context.Background(), req)
	if err != nil {
		t.Errorf("MGetSku error: %v", err)
		return
	}
	fmt.Printf("MGetSku response: %+v\n", resp)

	// 打印每个商品的详细信息
	for i, sku := range resp.Sku {
		fmt.Printf("Sku %d: %+v\n", i+1, sku)
	}
}

// 可选：添加测试数据的辅助函数
func TestInsertTestData(t *testing.T) {
	db := getTestDB()

	testSku := &database.GoodsSku{
		Sku:       "test_sku_001",
		GoodsID:   "test_goods_001",
		TagID:     "test_tag_001",
		Name:      "测试商品1",
		Price:     99.99,
		Spec:      "测试规格",
		ShowPic:   "pic1.jpg,pic2.jpg",
		DetailPic: "detail1.jpg,detail2.jpg",
		IsDeleted: 0,
	}

	result := db.Create(testSku)
	if result.Error != nil {
		t.Errorf("Insert test data error: %v", result.Error)
		return
	}
	fmt.Printf("Inserted test data: %+v\n", testSku)
}
