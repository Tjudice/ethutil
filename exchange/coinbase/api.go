package coinbase

// type CoinbaseAPI interface {
// 	MarketsAPI
// 	UserAPI
// 	TradeAPI
// 	WebsocketAPI
// }

// type MarketsAPI interface {
// 	GetMarkets(ctx context.Context) ([]*Market, error)
// 	GetMarket(ctx context.Context, marketId string) (*Market, error)
// 	GetMarketBookLevel1(ctx context.Context, marketId string) (*Orderbook, error)
// 	GetMarketBookLevel2(ctx context.Context, marketId string) (*Orderbook, error)
// 	GetMarketBookLevel3(ctx context.Context, marketId string) (*Orderbook, error)
// 	GetMarketCandles(ctx context.Context, marketId string) (*Candles, error)
// 	GetMarketTicker(ctx context.Context, marketId string) (*Ticker, error)
// 	GetMarketTrades(ctx context.Context, marketId string, limit int64) ([]*Trade, error)
// }

// type UserAPI interface {
// 	Accounts(ctx context.Context, limit *int32, cursor *string) ([]*Account, error)
// }

// type TradeAPI interface{}

// type WebsocketAPI interface {
// 	Subscribe(ctx context.Context, products []string, channel []any) (*ExchangeWebsocket, error)
// 	Unsubscribe(ctx context.Context, products []string, channel []any) error
// }
