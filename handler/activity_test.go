package handler

import (
	"context"
	"eshop_main/kitex_gen/eshop/home"
	"eshop_main/log"
	"testing"
)

func TestSubmitSeckillOrder(t *testing.T) {
	// 初始化测试请求
	req := &home.SeckillRequest{
		UserId:     "test_user_001",
		Sku:        "261311",
		ActivityId: "TEST_ACT_2023",
	}

	// 创建服务实例
	service := GoodsServiceImpl{}

	// 直接调用方法
	resp, err := service.SubmitSeckillOrder(context.Background(), req)

	// 简单输出结果
	t.Logf("测试结果: %+v 错误:%v", resp, err)
	select {}
}

func Test_cacheLotterySurplus(t *testing.T) {
	type args struct {
		ctx context.Context
		aId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx: context.Background(),
				aId: "testtt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cacheLotterySurplus(tt.args.ctx, tt.args.aId); (err != nil) != tt.wantErr {
				t.Errorf("cacheLotterySurplus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoodsServiceImpl_DrawLottery(t *testing.T) {
	type args struct {
		ctx context.Context
		req *home.LotteryRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *home.LotteryResponse
		wantErr bool
	}{
		{
			args: args{
				ctx: context.Background(),
				req: &home.LotteryRequest{
					UserId:     "ttttt",
					ActivityId: "testtt",
				},
			},
		},
	}
	for _, tt := range tests {
		for i := 0; i < 13; i++ {
			t.Run(tt.name, func(t *testing.T) {
				g := GoodsServiceImpl{}
				got, err := g.DrawLottery(tt.args.ctx, tt.args.req)
				if err != nil {
					log.Errorf("%d", err)
				}
				if got.HasWon {
					log.Infof("hasWon:%d, sku: %s, info: %s", got.HasWon, *got.Sku, got.Info)
				} else {
					log.Infof("hasWon:%d", got.HasWon)
				}
			})
		}
	}
}
