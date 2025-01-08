package main

import (
	"context"
	"eshop_main/database"
	"eshop_main/kitex_gen/eshop/home"
)

type GoodsServiceImpl struct{}

func (g GoodsServiceImpl) GetOneSku(ctx context.Context, sku string) (r *home.GetOneSkuResponse, err error) {
	goods, err := database.GetGoodsBySku(nil, sku)
	if err != nil {
		errStr := ""
		if err == database.ErrRecordNotFound {
			errStr = "商品不存在"
			return &home.GetOneSkuResponse{
				Code:   404,
				ErrStr: &errStr,
			}, nil
		}
		errStr = "服务器内部错误"
		return &home.GetOneSkuResponse{
			Code:   500,
			ErrStr: &errStr,
		}, err
	}

	return &home.GetOneSkuResponse{
		Code: 200,
		Sku:  database.ConvertToHomeGoodsSku(goods),
	}, nil
}

func (g GoodsServiceImpl) MGetSku(ctx context.Context, req *home.MGetSkuRequest) (r *home.MGetSkuResponse, err error) {
	skus, isEnd, err := database.GetGoodsList(nil, req.TagId, req.PageSize, req.PageNum)
	if err != nil {
		return &home.MGetSkuResponse{}, err
	}

	return &home.MGetSkuResponse{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		IsEnd:    isEnd,
		Sku:      skus,
	}, nil
}
