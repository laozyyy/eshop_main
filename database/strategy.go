package database

import (
	"errors"
	"eshop_main/log"
	"gorm.io/gorm"
	"time"
)

// LotteryStrategy 抽奖策略表
type LotteryStrategy struct {
	ID           uint      `gorm:"primaryKey;autoIncrement;comment:'ID，主键，自增'" json:"id"`
	StrategyID   string    `gorm:"unique;comment:'抽奖策略ID';size:255" json:"strategy_id"`
	StrategyName string    `gorm:"notNull;comment:'抽奖策略名称，唯一标识一个策略';size:255" json:"strategy_name"`
	StartTime    time.Time `gorm:"notNull;comment:'抽奖策略生效的开始时间'" json:"start_time"`
	EndTime      time.Time `gorm:"notNull;comment:'抽奖策略生效的结束时间'" json:"end_time"`
	Score        int       `gorm:"comment:'抽奖使用积分'" json:"score"`
	Surplus      int       `gorm:"comment:'抽奖总库存'" json:"surplus"`
	Status       int8      `gorm:"default:1;comment:'抽奖策略状态：1-启用，2-停用'" json:"status"`
	CreateTime   time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime   time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

// LotteryStrategyAward 抽奖策略奖品表
type LotteryStrategyAward struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;comment:'ID，主键，自增'" json:"id"`
	StrategyID string    `gorm:"unique;comment:'抽奖策略ID';size:255" json:"strategy_id"`
	SKU        string    `gorm:"comment:'sku，为0代表未中奖';size:255" json:"sku"`
	Surplus    int       `gorm:"comment:'sku库存'" json:"surplus"`
	Weight     int       `gorm:"comment:'权重'" json:"weight"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"create_time"`
	UpdateTime time.Time `gorm:"default:CURRENT_TIMESTAMP;update:CURRENT_TIMESTAMP;comment:'更新时间'" json:"update_time"`
}

func GetStrategy(db *gorm.DB, aId string) (*LotteryStrategy, error) {
	db = getDBInstance(db)
	var res []*LotteryStrategy
	err := db.Table("lottery_strategy").
		Where("strategy_id = ?", aId).
		First(&res).Error
	if err != nil {
		log.Errorf("error: %d", err)
		return nil, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return nil, errors.New("策略不存在")
}

func GetStrategyAward(db *gorm.DB, aId string) ([]*LotteryStrategyAward, error) {
	db = getDBInstance(db)
	var res []*LotteryStrategyAward
	err := db.Table("lottery_strategy_award").
		Where("strategy_id = ?", aId).
		Find(&res).Error
	if err != nil {
		log.Errorf("error: %d", err)
		return nil, err
	}
	if len(res) > 0 {
		return res, nil
	}
	return nil, errors.New("策略奖品没有内容")
}
