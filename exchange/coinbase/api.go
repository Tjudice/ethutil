package coinbase

import "context"

type CoinbaseAPI interface {
	MarketsAPI
	UserAPI
	TradeAPI
}

type MarketsAPI interface {
	GetMarkets(ctx context.Context) ([]*Market, error)
}

type UserAPI interface {
	Accounts(ctx context.Context, limit *int32, cursor *string) ([]*Account, error)
}

type TradeAPI interface{}
