// Code generated by Kitex v0.12.0. DO NOT EDIT.

package goodsservice

import (
	"context"
	home "eshop_main/kitex_gen/eshop/home"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	GetOneSku(ctx context.Context, sku string, callOptions ...callopt.Option) (r *home.GetOneSkuResponse, err error)
	GetRandomSku(ctx context.Context, req *home.PageRequest, callOptions ...callopt.Option) (r *home.PageResponse, err error)
	MGetSku(ctx context.Context, sku *home.MGetSkuRequest, callOptions ...callopt.Option) (r *home.PageResponse, err error)
	GetPrice(ctx context.Context, req *home.GetPriceRequest, callOptions ...callopt.Option) (r string, err error)
	SearchGoods(ctx context.Context, req *home.SearchRequest, callOptions ...callopt.Option) (r *home.PageResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfoForClient(), options...)
	if err != nil {
		return nil, err
	}
	return &kGoodsServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kGoodsServiceClient struct {
	*kClient
}

func (p *kGoodsServiceClient) GetOneSku(ctx context.Context, sku string, callOptions ...callopt.Option) (r *home.GetOneSkuResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetOneSku(ctx, sku)
}

func (p *kGoodsServiceClient) GetRandomSku(ctx context.Context, req *home.PageRequest, callOptions ...callopt.Option) (r *home.PageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetRandomSku(ctx, req)
}

func (p *kGoodsServiceClient) MGetSku(ctx context.Context, sku *home.MGetSkuRequest, callOptions ...callopt.Option) (r *home.PageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MGetSku(ctx, sku)
}

func (p *kGoodsServiceClient) GetPrice(ctx context.Context, req *home.GetPriceRequest, callOptions ...callopt.Option) (r string, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetPrice(ctx, req)
}

func (p *kGoodsServiceClient) SearchGoods(ctx context.Context, req *home.SearchRequest, callOptions ...callopt.Option) (r *home.PageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SearchGoods(ctx, req)
}
