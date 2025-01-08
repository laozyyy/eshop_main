package main

import (
	"context"
	"eshop_main/database"
	"eshop_main/kitex_gen/eshop/home"
	"eshop_main/log"
)

type GoodsServiceImpl struct{}

func (g GoodsServiceImpl) GetOneSku(ctx context.Context, sku string) (r *home.GetOneSkuResponse, err error) {
	log.Infof("请求获取 SKU: %s", sku)
	goods, err := database.GetGoodsBySku(nil, sku)
	if err != nil {
		errStr := ""
		if err == database.ErrRecordNotFound {
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

func (g GoodsServiceImpl) MGetSku(ctx context.Context, req *home.MGetSkuRequest) (r *home.MGetSkuResponse, err error) {
	log.Infof("请求批量获取 SKU, TagID: %s, PageSize: %d, PageNum: %d", req.TagId, req.PageSize, req.PageNum)
	skus, isEnd, err := database.GetGoodsList(nil, req.TagId, req.PageSize, req.PageNum)
	if err != nil {
		log.Errorf("批量获取 SKU 时发生错误: %v", err)
		return &home.MGetSkuResponse{}, err
	}

	log.Infof("成功获取 SKU 列表, 是否结束: %v", isEnd)
	return &home.MGetSkuResponse{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		IsEnd:    isEnd,
		Sku:      skus,
	}, nil
}