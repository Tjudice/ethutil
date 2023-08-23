package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tjudice/util/go/lambda"
	"github.com/tjudice/util/go/network/http_helpers"
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
	return http_helpers.GetJSONFn[*BestBidAsks](ctx, c.cl, ADVANCED_TRADE_BEST_BID_ASK_URL, nil, func(r *http.Request) {
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).Add("product_ids", lambda.SliceToAny(productIds)...).Encode()
	})
}

type OrderbookWrapperAdvancedTrade struct {
	Pricebook *AdvancedTradeOrderbook `json:"pricebook"`
}

type AdvancedTradeOrderbook struct {
	ProductId string    `json:"product_id"`
	Bids      []*Tick   `json:"bids"`
	Asks      []*Tick   `json:"asks"`
	Time      time.Time `json:"time"`
}

const ADVANCED_TRADE_ORDERBOOK_URL = "https://api.coinbase.com/api/v3/brokerage/product_book"

func (c *AdvancedTradeClient) GetOrderbook(ctx context.Context, productId string, limit int) (*AdvancedTradeOrderbook, error) {
	wrapper, err := http_helpers.GetJSONFn[*OrderbookWrapperAdvancedTrade](ctx, c.cl, ADVANCED_TRADE_ORDERBOOK_URL, nil, func(r *http.Request) {
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).Add("product_id", productId).AddIfNotDefault("limit", limit, 0).Encode()
	})
	if err != nil {
		return nil, err
	}
	return wrapper.Pricebook, nil
}

type GetMarketParams struct {
	Limit              int
	Offset             int
	ProductType        ProductType
	ProductIds         []string
	ContractExpiryType ContractExpiryType
}

type AdvancedTradeMarket struct {
	ProductId                 string      `json:"product_id"`
	Price                     float64     `json:"price,string"`
	PricePercentageChange24H  float64     `json:"price_percentage_change_24h,string"`
	Volume24H                 float64     `json:"volume_24h,string"`
	VolumePercentageChange24H float64     `json:"volume_percentage_change_24h,string"`
	BaseIncrement             float64     `json:"base_increment,string"`
	QuoteIncrement            float64     `json:"quote_increment,string"`
	QuoteMinSize              float64     `json:"quote_min_size,string"`
	QuoteMaxSize              float64     `json:"quote_max_size,string"`
	BaseMinSize               float64     `json:"base_min_size,string"`
	BaseMaxSize               float64     `json:"base_max_size,string"`
	BaseName                  string      `json:"base_name"`
	QuoteName                 string      `json:"quote_name"`
	Watched                   bool        `json:"watched"`
	IsDisabled                bool        `json:"is_disabled"`
	New                       bool        `json:"new"`
	Status                    string      `json:"status"`
	CancelOnly                bool        `json:"cancel_only"`
	LimitOnly                 bool        `json:"limit_only"`
	PostOnly                  bool        `json:"post_only"`
	TradingDisabled           bool        `json:"trading_disabled"`
	AuctionMode               bool        `json:"auction_mode"`
	ProductType               string      `json:"product_type"`
	QuoteCurrencyId           string      `json:"quote_currency_id"`
	BaseCurrencyId            string      `json:"base_currency_id"`
	FcmTradingSessionDetails  interface{} `json:"fcm_trading_session_details"`
	MidMarketPrice            string      `json:"mid_market_price"`
	AliasTo                   []string    `json:"alias_to"`
	BaseDisplaySymbol         string      `json:"base_display_symbol"`
	QuoteDisplaySymbol        string      `json:"quote_display_symbol"`
	ViewOnly                  bool        `json:"view_only"`
	PriceIncrement            float64     `json:"price_increment,string"`
}

type AdvancedTradeMarkets struct {
	Products []json.RawMessage `json:"products"`
}

const ADVANCED_TRADE_GET_MARKETS_URL = "https://api.coinbase.com/api/v3/brokerage/products"

func (c *AdvancedTradeClient) GetMarkets(ctx context.Context, params *GetMarketParams) (json.RawMessage, error) {
	return http_helpers.GetJSONFn[json.RawMessage](ctx, c.cl, ADVANCED_TRADE_GET_MARKETS_URL, nil, func(r *http.Request) {
		if params == nil {
			return
		}
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).
			AddIfNotDefault("limit", params.Limit, 0).
			AddIfNotDefault("offset", params.Offset, 0).
			AddIfNotDefault("product_type", string(params.ProductType), "").
			Add("product_ids", lambda.SliceToAny(params.ProductIds)...).
			AddIfNotDefault("contract_expiry_type", string(params.ContractExpiryType), "").Encode()
	})
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
