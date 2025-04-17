package database

import (
	"eshop_main/log"
	"gorm.io/gorm"
)

type GoodsSkuSurplus struct {
	ID      int64  `gorm:"primaryKey;column:id"`
	Sku     string `gorm:"column:sku"`
	Surplus string `gorm:"column:surplus"`
}

func UpdateSurplus(db *gorm.DB, sku string, surplus string) error {
	db = getDBInstance(db)
	err := db.Table("goods_sku_surplus").
		Where("sku = ?", sku).
		Update("surplus", surplus).Error
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	return nil
}

func GetSurplus(db *gorm.DB, sku string) (string, error) {
	db = getDBInstance(db)
	var res []*GoodsSkuSurplus
	err := db.Table("goods_sku_surplus").
		Where("sku = ?", sku).
		Find(&res).Error
	if err != nil {
		log.Errorf("error: %v", err)
		return "", err
	}
	if len(res) > 0 {
		return res[0].Surplus, nil
	}
	return "", nil
}
