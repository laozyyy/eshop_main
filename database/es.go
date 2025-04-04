package database

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

var ESClient *elastic.Client

const goodsMapping = `{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
    },
    "mappings": {
        "properties": {
            "sku":        { "type": "keyword" },
            "name":       { "type": "text", "analyzer": "ik_max_word" },
            "price":      { "type": "double" },
            "spec":       { "type": "text" },
            "tag_id":     { "type": "keyword" },
            "is_deleted": { "type": "integer" }
        }
    }
}`

func init() {
	_ = InitES()
	//if err != nil {
	//	panic(err)
	//}
}

func InitES() error {
	client, err := elastic.NewClient(
		elastic.SetURL("http://117.72.72.114:19200"),
		elastic.SetBasicAuth("elastic", "123456"), // 添加认证信息
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return fmt.Errorf("ES连接失败: %w", err)
	}
	ESClient = client

	// 创建商品索引
	ctx := context.Background()
	exists, err := client.IndexExists("goods_index").Do(ctx)
	if err != nil {
		return fmt.Errorf("索引检查失败: %w", err)
	}

	if !exists {
		createIndex, err := client.CreateIndex("goods_index").
			BodyString(goodsMapping).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("索引创建失败: %w", err)
		}
		if !createIndex.Acknowledged {
			return fmt.Errorf("索引创建未被确认")
		}
	}
	return nil
}
