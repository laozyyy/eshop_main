package handler

import (
	"context"
	"errors"
	"eshop_main/database"
	"eshop_main/kitex_gen/eshop/home"
	"eshop_main/log"
	"strconv"
)

type GoodsServiceImpl struct{}

func (g GoodsServiceImpl) GetPrice(ctx context.Context, req *home.GetPriceRequest) (r string, err error) {
	log.Infof("请求获取 s: %s", req.Sku)
	goods, err := database.GetGoodsBySku(nil, req.Sku)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			log.Errorf("SKU %s 不存在: %v", req.Sku, err)
			return "", nil
		}
		log.Errorf("获取 SKU %s 时发生错误: %v", req.Sku, err)
		return "", nil
	}

	log.Infof("成功获取 SKU: %+v", goods)
	ret := strconv.FormatFloat(goods.Price, 'f', 2, 64)
	return ret, nil
}

func (g GoodsServiceImpl) GetOneSku(ctx context.Context, sku string) (r *home.GetOneSkuResponse, err error) {
	log.Infof("请求获取 SKU: %s", sku)
	goods, err := database.GetGoodsBySku(nil, sku)
	if err != nil {
		errStr := ""
		if errors.Is(err, database.ErrRecordNotFound) {
			errStr = "商品不存在"
			log.Errorf("SKU %s 不存在: %v", sku, err)
			return &home.GetOneSkuResponse{
				Code:   404,
				ErrStr: &errStr,
			}, nil
		}
		errStr = "服务器内部错误"
		log.Errorf("获取 SKU %s 时发生错误: %v", sku, err)
		return &home.GetOneSkuResponse{
			Code:   500,
			ErrStr: &errStr,
		}, err
	}

	log.Infof("成功获取 SKU: %+v", goods)
	return &home.GetOneSkuResponse{
		Code: 200,
		Sku:  database.ConvertToHomeGoodsSku(goods),
	}, nil
}

func (g GoodsServiceImpl) MGetSku(ctx context.Context, req *home.MGetSkuRequest) (r *home.PageResponse, err error) {
	log.Infof("请求批量获取 SKU, TagID: %s, PageSize: %d, PageNum: %d", req.TagId, req.PageSize, req.PageNum)
	skus, isEnd, err := database.GetGoodsList(nil, req.TagId, req.PageSize, req.PageNum)
	if err != nil {
		log.Errorf("批量获取 SKU 时发生错误: %v", err)
		return &home.PageResponse{}, err
	}

	log.Infof("成功获取 SKU 列表, 是否结束: %v", isEnd)
	return &home.PageResponse{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		IsEnd:    isEnd,
		Sku:      skus,
	}, nil
}

func (g GoodsServiceImpl) GetRandomSku(ctx context.Context, req *home.PageRequest) (r *home.PageResponse, err error) {
	log.Infof("请求批量随机获取 SKU, PageSize: %d, PageNum: %d", req.PageSize, req.PageNum)
	skus, _, err := database.GetRandomGoodsList(nil, req.PageSize)
	if err != nil {
		log.Errorf("批量获取 SKU 时发生错误: %v", err)
		return &home.PageResponse{}, err
	}

	return &home.PageResponse{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		IsEnd:    false,
		Sku:      skus,
	}, nil
}

func (g GoodsServiceImpl) SearchGoods(ctx context.Context, req *home.SearchRequest) (r *home.PageResponse, err error) {
	log.Infof("商品搜索请求 关键词:%s 分页:%d/%d", req.Keyword, req.PageNum, req.PageSize)

	skus, isEnd, err := database.SearchGoodsByName(req.Keyword, req.PageSize, req.PageNum)
	if err != nil {
		log.Errorf("商品搜索失败 关键词:%s 错误:%v", req.Keyword, err)
		return &home.PageResponse{}, err
	}

	return &home.PageResponse{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		IsEnd:    isEnd,
		Sku:      skus,
	}, nil
}
