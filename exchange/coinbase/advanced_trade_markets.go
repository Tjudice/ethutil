package coinbase

import "context"

type BidAsk struct{}

const ADVANCED_TRADE_BEST_BID_ASK_URL = "https://api.coinbase.com/api/v3/brokerage/best_bid_ask"

func (c *AdvancedTradeClient) GetBestBidAsk(ctx context.Context, productIds []string) ([]*BidAsk, error) {
	panic("not implemented")
}

type AdvancedTradeOrderbook struct{}

const ADVANCED_TRADE_ORDERBOOK_URL = "https://api.coinbase.com/api/v3/brokerage/product_book"

func (c *AdvancedTradeClient) GetOrderbook(ctx context.Context, limit int) (*AdvancedTradeOrderbook, error) {
	panic("not implemented")
}

type GetMarketParams struct {
	Limit              int
	Offset             int
	ProductType        ProductType
	ProductIds         []string
	ContractExpiryType ContractExpiryType
}

type AdvancedTradeMarket struct{}

type AdvancedTradeMarkets struct{}

const ADVANCED_TRADE_GET_MARKETS_URL = "https://api.coinbase.com/api/v3/brokerage/products"

func (c *AdvancedTradeClient) GetMarkets(ctx context.Context, params *GetMarketParams) (*AdvancedTradeMarkets, error) {
	panic("not implemented")
}

const ADVANCED_TRADE_GET_MARKET_URL = "https://api.coinbase.com/api/v3/brokerage/products/%s"

func (c *AdvancedTradeClient) GetMarket(ctx context.Context, marketId string) (*AdvancedTradeMarket, error) {
	panic("not implemented")
}

type AdvancedTradeCandle struct{}

type AdvandedTradeCandles struct{}

const ADVANCED_TRADE_CANDLES_URL = "https://api.coinbase.com/api/v3/brokerage/products/%s/candles"

func (c *AdvancedTradeClient) GetCandles(ctx context.Context, marketId string, granularity, start, stop int64) (*AdvandedTradeCandles, error) {
	panic("not implemented")
}

type AdvancedTradeMarketTrade struct{}

type AdvancedTradeMarketTrades struct{}

const ADVANCED_TRADE_MARKET_TRADES_URL = "https://api.coinbase.com/api/v3/brokerage/products/%s/ticker"

func (c *AdvancedTradeClient) GetMarketTrades(ctx context.Context, marketId string, limit int) (*AdvancedTradeMarketTrades, error) {
	panic("not implemented")
}
