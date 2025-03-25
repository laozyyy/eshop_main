package database

import (
	"context"
	"encoding/json"
	"eshop_main/kitex_gen/eshop/home"
	"eshop_main/log"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

type GoodsSku struct {
	ID         int64     `gorm:"primaryKey;column:id"`
	Sku        string    `gorm:"column:sku"`
	GoodsID    string    `gorm:"column:goods_id"`
	TagID      string    `gorm:"column:tag_id"`
	Name       string    `gorm:"column:name"`
	Price      float64   `gorm:"column:price"`
	Spec       string    `gorm:"column:spec"`
	ShowPic    string    `gorm:"column:show_pic"`
	DetailPic  string    `gorm:"column:detail_pic"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
	IsDeleted  int32     `gorm:"column:is_deleted"`
}

func (g *GoodsSku) TableName() string {
	return "goods_sku"
}

func getDBInstance(db *gorm.DB) *gorm.DB {
	if db == nil {
		if DB == nil {
			Init() // 初始化全局 DB
		}
		return DB
	}
	return db
}

func GetGoodsBySku(db *gorm.DB, sku string) (*GoodsSku, error) {
	db = getDBInstance(db)
	var goods GoodsSku
	result := db.Where("sku = ? AND is_deleted = 0", sku).First(&goods)
	if result.Error != nil {
		return nil, result.Error
	}
	return &goods, nil
}

func GetGoodsBySkus(db *gorm.DB, skus []string) ([]*GoodsSku, error) {
	db = getDBInstance(db)
	var goodsList []*GoodsSku
	result := db.Where("sku IN ? AND is_deleted = 0", skus).Find(&goodsList)
	if result.Error != nil {
		return nil, result.Error
	}
	return goodsList, nil
}

func GetGoodsList(db *gorm.DB, tagID string, pageSize, pageNum int32) ([]*home.Sku, bool, error) {
	db = getDBInstance(db)
	var goodsList []*GoodsSku
	query := db.Where("is_deleted = 0")

	if tagID != "" {
		query = query.Where("tag_id = ?", tagID)
	}

	offset := (pageNum - 1) * pageSize

	err := query.Offset(int(offset)).
		Limit(int(pageSize + 1)). // 多查询一条用于判断是否结束
		Find(&goodsList).Error

	if err != nil {
		log.Errorf("error: %v", err)
		return nil, false, err
	}

	isEnd := true
	if len(goodsList) > int(pageSize) {
		isEnd = false
		goodsList = goodsList[:pageSize] // 去掉多查询的一条
	}

	var skus []*home.Sku
	for _, goods := range goodsList {
		skus = append(skus, ConvertToHomeGoodsSku(goods))
	}

	return skus, isEnd, nil
}

func GetRandomGoodsList(db *gorm.DB, pageSize int32) ([]*home.Sku, bool, error) {
	db = getDBInstance(db)
	var goodsList []*GoodsSku
	var maxID int
	err := db.Model(&GoodsSku{}).
		Select("COALESCE(MAX(id), 0)").
		Scan(&maxID).
		Error
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, false, err
	}

	var randomIDs []int
	for i := 0; i < int(pageSize); i++ {
		randomID := rand.Intn(maxID) + 1
		randomIDs = append(randomIDs, randomID)
	}

	// 查询ID >= randomID 的第一条记录（避免空洞ID）
	// todo 可能查到的数量不够pagesize
	err = db.Model(&GoodsSku{}).
		Where("id >= ?", randomIDs[0]).
		Limit(int(pageSize)).
		Find(&goodsList).
		Error
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, false, err
	}

	var skus []*home.Sku
	for _, goods := range goodsList {
		skus = append(skus, ConvertToHomeGoodsSku(goods))
	}

	return skus, false, nil
}

func ConvertToHomeGoodsSku(goods *GoodsSku) *home.Sku {
	// 将字符串转换为字符串数组
	showPics := strings.Split(goods.ShowPic, ",")
	detailPics := strings.Split(goods.DetailPic, ",")

	return &home.Sku{
		Sku:        goods.Sku,
		GoodsId:    goods.GoodsID,
		TagId:      goods.TagID,
		Name:       goods.Name,
		Price:      int32(goods.Price * 100), // 转换为分为单位的整数
		Spec:       goods.Spec,
		ShowPic:    showPics,
		DetailPic:  detailPics,
		SellerName: "", // 需要补充
	}
}

var ErrRecordNotFound = gorm.ErrRecordNotFound

// 修改现有的SearchGoodsByName实现
func SearchGoodsByName(keyword string, pageSize, pageNum int32) ([]*home.Sku, bool, error) {
	// 使用ES客户端进行搜索
	searchResult, err := ESClient.Search().
		Index("goods_index").
		Query(elastic.NewMatchQuery("name", keyword)).
		From(int((pageNum - 1) * pageSize)).
		Size(int(pageSize)).
		Do(context.Background())

	if err != nil {
		log.Errorf("ES查询失败 关键词:%s 错误:%v", keyword, err)
		return nil, false, err
	}

	// 解析ES结果
	var skus []*home.Sku
	for _, hit := range searchResult.Hits.Hits {
		var goods GoodsSku
		if err := json.Unmarshal(hit.Source, &goods); err != nil {
			log.Errorf("ES数据解析失败 ID:%s 错误:%v", hit.Id, err)
			continue
		}
		skus = append(skus, ConvertToHomeGoodsSku(&goods))
	}

	// 计算是否结束
	total := searchResult.TotalHits()
	currentCount := int64(pageNum * pageSize)
	isEnd := currentCount >= total

	return skus, isEnd, nil
}

// 在商品创建/更新时添加ES同步
func CreateGoods(db *gorm.DB, goods *GoodsSku) error {
	// 原有数据库操作
	tx := db.Begin()
	if err := tx.Create(goods).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 新增ES同步
	_, err := ESClient.Index().
		Index("goods_index").
		Id(goods.Sku).
		BodyJson(goods).
		Refresh("wait_for").
		Do(context.Background())

	if err != nil {
		tx.Rollback()
		log.Errorf("ES同步失败 SKU:%s 错误:%v", goods.Sku, err)
		return err
	}

	tx.Commit()
	return nil
}
