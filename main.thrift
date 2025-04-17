namespace go eshop.home

service goodsService {
    GetOneSkuResponse GetOneSku(1: string sku)
    PageResponse GetRandomSku(1: PageRequest req)
    PageResponse MGetSku(1: MGetSkuRequest sku)

    PageResponse SearchGoods(1: SearchRequest req)

    string GetPrice(1: GetPriceRequest req)

    SeckillResponse SubmitSeckillOrder(1: SeckillRequest req)
    LotteryResponse DrawLottery(1: LotteryRequest req)
}
struct GetPriceRequest {
    1: string sku
}
struct PageRequest {
    1: i32 pageSize
    2: i32 pageNum
}

struct PageResponse {
    1: i32 pageSize
    2: i32 pageNum
    3: bool isEnd
    4: list<Sku> sku
}
struct MGetSkuRequest {
    1: i32 pageSize
    2: i32 pageNum
    3: string tagId
}
struct GetOneSkuResponse {
    1: optional Sku sku
    2: i64 code
    3: optional string errStr
}

struct Sku {
    1: string sku,
    2: string goodsId,
    3: string tagId,
    4: string name,
    5: i32 price,
    6: string spec,
    7: list<string> showPic,
    8: list<string> detailPic,
    9: string sellerName
}

struct SearchRequest {
    1: string keyword
    2: i32 pageSize
    3: i32 pageNum
}


// 新增秒杀相关结构体
struct SeckillRequest {
    1: string userId,
    2: string sku,
    3: string activityId
}

struct SeckillResponse {
    1: string info
    2: string orderId
}

// 新增抽奖相关结构体
struct LotteryRequest {
    1: string userId,
    2: string activityId
}

struct LotteryResponse {
    1: bool hasWon, // 是否中奖
    2: optional string sku,
    3: string info // 是否成功
    4: string orderId
}

//kitex -module eshop_main main.thrift
