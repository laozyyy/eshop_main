namespace go eshop.home

service goodsService {
    GetOneSkuResponse GetOneSku(1: string sku)
    MGetSkuResponse MGetSku(1: MGetSkuRequest sku)
}

struct MGetSkuResponse {
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
    1: string sellerId,
    2: string goodsId,
    3: string tagId,
    4: string name,
    5: i32 price,
    6: string spec,
    7: list<string> showPic,
    8: list<string> detailPic,
    9: string sellerName
}