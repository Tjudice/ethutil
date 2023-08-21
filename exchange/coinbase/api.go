package coinbase

import "context"

type CoinbaseAPI interface {
	MarketsAPI
	UserAPI
	TradeAPI
}

type MarketsAPI interface {
	GetMarkets(ctx context.Context) ([]*Market, error)
	GetMarket(ctx context.Context, marketId string) (*Market, error)
	GetMarketBookLevel1(ctx context.Context, marketId string) (*Orderbook, error)
	GetMarketBookLevel2(ctx context.Context, marketId string) (*Orderbook, error)
	GetMarketBookLevel3(ctx context.Context, marketId string) (*Orderbook, error)
}

type UserAPI interface {
	Accounts(ctx context.Context, limit *int32, cursor *string) ([]*Account, error)
}

type TradeAPI interface{}
