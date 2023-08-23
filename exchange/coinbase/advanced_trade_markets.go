package coinbase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tjudice/util/go/clients/jsonhttp"
)

type BestBidAsks struct {
	PriceBooks []*BidAsk `json:"pricebooks"`
}

type BidAsk struct {
	ProductId string    `json:"product_id"`
	Bids      []*Tick   `json:"bids"`
	Asks      []*Tick   `json:"asks"`
	Time      time.Time `json:"time"`
}

type Tick struct {
	Price float64 `json:"price,string"`
	Size  float64 `json:"size,string"`
}

const ADVANCED_TRADE_BEST_BID_ASK_URL = "https://api.coinbase.com/api/v3/brokerage/best_bid_ask"

func (c *AdvancedTradeClient) GetBestBidAsk(ctx context.Context, productIds []string) (*BestBidAsks, error) {
	if len(productIds) == 0 {
		return nil, fmt.Errorf("GetBestBidAsk: must provide at least 1 product id")
	}
	return jsonhttp.Get[*BestBidAsks](ctx, c.cl, ADVANCED_TRADE_BEST_BID_ASK_URL+"?product_ids="+strings.Join(productIds, "&product_ids="), nil)
}

type AdvancedTradeOrderbook struct{}

const ADVANCED_TRADE_ORDERBOOK_URL = "https://api.coinbase.com/api/v3/brokerage/product_book"

func (c *AdvancedTradeClient) GetOrderbook(ctx context.Context, productId string, limit int) (*AdvancedTradeOrderbook, error) {

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
