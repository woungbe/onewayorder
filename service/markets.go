package service

type ExchnageInfo struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int64         `json:"serverTime"`
	RateLimits      []RateLimit   `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []SymbolData  `json:"symbols"`
}

type ExchangeTickInfo struct {
	// 4개로 나누자
	CoinDecimalSize  string // 코인 최소 소수점  ps) 보통 최소 소수점이랑 tick size 랑 동일하게 사용해서 헷갈림
	CoinTickSize     string // 코인 tick size
	PriceDecimalSize string // 가격 소수점 - 가격 기준 소수점이 얼마인가
	PriceTickSize    string // 가격 tick 사이즈 - 가격 기준 한틱이 얼마인가.

	// 주문시 최소 수량
	MinQuantity string // 주문시 최소 거래 수량
	Minnotional string // 주문시 최소 거래 금액

	// 주문시
	MarketMaxQuantity string // 시장가 주문시 최대 거래 수량
	MarketMaxPrice    string // 시장가 주문시 최대 거래 금액
}

type FuturesMarket struct {
	Markets        map[string]ExchnageInfo
	MapLotTickInfo map[string]*ExchangeTickInfo
}
