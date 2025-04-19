package handler

import (
	"context"
	"errors"
	"eshop_main/cache"
	"eshop_main/database"
	"eshop_main/kitex_gen/eshop/home"
	"eshop_main/log"
	"eshop_main/mq"
	"fmt"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// 秒杀订单提交
func (g GoodsServiceImpl) SubmitSeckillOrder(ctx context.Context, req *home.SeckillRequest) (*home.SeckillResponse, error) {
	log.Infof("秒杀请求 用户:%s 商品:%s", req.UserId, req.Sku)
	surplusKey := getSecKillSurplusKey(req.ActivityId)
	err := cacheSecKillSurplus(ctx, req.ActivityId, req.Sku)
	if err != nil {
		log.Errorf("写入缓存失败 error: %d", err)
		return &home.SeckillResponse{
			Info: "内部错误",
		}, nil
	}
	result, err := cache.Client.Decr(ctx, surplusKey).Result()
	if err != nil {
		log.Errorf("库存:%s 扣减错误: %v", req.Sku, err)
		return &home.SeckillResponse{
			Info: "库存扣减错误",
		}, nil
	}
	if result < 0 {
		return &home.SeckillResponse{
			Info: "已无库存",
		}, nil
	}

	// 异步执行库存扣减
	go mq.SendSurplusDecr(req.Sku, req.ActivityId)

	// todo 调用订单服务，生成订单

	newUUID, _ := uuid.NewUUID()
	return &home.SeckillResponse{
		Info:    "success",
		OrderId: newUUID.String(), // 需实现订单号生成逻辑
	}, nil
}

func cacheSecKillSurplus(ctx context.Context, aId string, sku string) error {
	surplusKey := getSecKillSurplusKey(aId)
	exist, _ := cache.Client.Exists(ctx, surplusKey).Result()
	if exist == 0 {
		surplus, err := database.GetSurplus(nil, sku)
		if err != nil {
			if err != nil {
				log.Errorf("error: %d", err)
			}
			return err
		}
		if surplus == "" {
			return errors.New(fmt.Sprintf("库存异常，sku：%s, aid: %s", sku, aId))
		}
		cache.Client.Set(ctx, surplusKey, surplus, time.Hour*12)
	}
	return nil
}

// 抽奖功能
func (g GoodsServiceImpl) DrawLottery(ctx context.Context, req *home.LotteryRequest) (*home.LotteryResponse, error) {
	log.Infof("抽奖请求 用户:%s 活动:%s", req.UserId, req.ActivityId)
	surplusKey := getLotterySurplusKey(req.ActivityId)
	err := cacheLotterySurplus(ctx, req.ActivityId)
	if err != nil {
		log.Errorf("写入缓存失败")
		return &home.LotteryResponse{
			Info: "内部错误",
		}, nil
	}
	err = cache.Client.Decr(ctx, surplusKey).Err()
	if err != nil {
		log.Errorf("抽奖库存:%s 扣减错误: %v", req.ActivityId, err)
		return &home.LotteryResponse{
			Info: "库存扣减错误",
		}, nil
	}
	awardId, err := raffle(ctx, req.ActivityId)
	if err != nil {
		log.Errorf("抽奖:%s 抽奖错误: %v", req.ActivityId, err)
		return &home.LotteryResponse{
			Info: "抽奖错误",
		}, nil
	}
	if awardId == "0" {
		return &home.LotteryResponse{
			Info:   "success",
			HasWon: false,
		}, nil
	}
	// 异步执行库存扣减
	go mq.SendLotterySurplusDecr(req.ActivityId)

	// todo 调用订单服务，生成订单

	newUUID, _ := uuid.NewUUID()

	return &home.LotteryResponse{
		HasWon:  true,
		Sku:     &awardId,
		Info:    "success",
		OrderId: newUUID.String(), // 需实现订单号生成逻辑
	}, nil
}

func cacheLotterySurplus(ctx context.Context, aId string) error {
	hashKey := getLotteryHashKey(aId)
	lotterySurplusKey := getLotterySurplusKey(aId)
	exist, _ := cache.Client.Exists(ctx, hashKey).Result()
	if exist == 1 {
		return nil
	}
	strategy, err := database.GetStrategy(nil, aId)
	if err != nil {
		log.Errorf("error : %d", err)
		return err
	}
	cache.Client.Set(ctx, lotterySurplusKey, strategy.Surplus, time.Hour*12)
	awards, err := database.GetStrategyAward(nil, aId)
	totalWeight := 0
	for _, award := range awards {
		if award.SKU != "0" {
			awardKey := getLotteryAwardSurplusKey(award.StrategyID, award.SKU)
			cache.Client.Set(ctx, awardKey, award.Surplus, time.Hour*12)
		}
		totalWeight += award.Weight
	}
	// 分配区间
	intervalMap := make(map[string][]int)
	currentStart := 1
	for _, award := range awards {
		// 计算当前商品的区间大小
		intervalSize := int(math.Round(float64(award.Weight) / float64(totalWeight) * 100))
		intervalEnd := currentStart + intervalSize - 1

		// 确保区间不超过100
		if intervalEnd > 100 {
			intervalEnd = 100
		}

		// 存储区间
		intervalMap[award.SKU] = []int{currentStart, intervalEnd}

		// 更新下一个区间的起始点
		currentStart = intervalEnd + 1

		// 如果已经分配完100，退出循环
		if currentStart > 100 {
			break
		}
	}

	// 打印区间分配结果
	for sku, interval := range intervalMap {
		log.Infof("SKU: %s, Interval: %d-%d", sku, interval[0], interval[1])
	}

	toCache := make(map[string]string)
	// 存储到Redis的Hash结构
	for sku, interval := range intervalMap {
		// 将区间范围存储为一个字符串
		for i := interval[0]; i <= interval[1]; i++ {
			toCache[strconv.Itoa(i)] = sku
		}

	}
	_, err = cache.Client.HSet(ctx, hashKey, toCache).Result()
	if err != nil {
		log.Errorf("error: %d", err)
	}
	return nil
}

func raffle(ctx context.Context, aId string) (string, error) {
	// 设置随机数种子，使用当前时间作为种子
	rand.Seed(time.Now().UnixNano())
	key := getLotteryHashKey(aId)

	// 生成1到1000的随机数
	randomNumber := rand.Intn(100) + 1
	sku, err := cache.Client.HGet(ctx, key, strconv.Itoa(randomNumber)).Result()
	if sku == "0" {
		return "0", nil
	}
	//扣减库存
	awardKey := getLotteryAwardSurplusKey(aId, sku)
	i, err := cache.Client.Decr(ctx, awardKey).Result()
	if i < 0 {
		return "0", nil
	}
	return sku, err
}

func getLotteryHashKey(id string) string {
	return fmt.Sprintf("lottery_table:%s", id)
}
func getLotteryAwardSurplusKey(aId string, sku string) string {
	return fmt.Sprintf("lottery_award_surplus:%s_%s", aId, sku)
}
func getLotterySurplusKey(id string) string {
	return fmt.Sprintf("lottery_surplus:%s", id)
}
func getSecKillSurplusKey(id string) string {
	return fmt.Sprintf("seckill_surplus:%s", id)
}
