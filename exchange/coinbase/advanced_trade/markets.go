package advanced_trade

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tjudice/util/go/lambda"
	"github.com/tjudice/util/go/network/http_helpers"
	"github.com/valyala/fastjson/fastfloat"
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

func (c *Client) GetBestBidAsk(ctx context.Context, productIds []string) (*BestBidAsks, error) {
	if len(productIds) == 0 {
		return nil, fmt.Errorf("GetBestBidAsk: must provide at least 1 product id")
	}
	return http_helpers.GetJSONFn[*BestBidAsks](ctx, c.cl, ADVANCED_TRADE_BEST_BID_ASK_URL, nil, func(r *http.Request) {
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).Add("product_ids", lambda.SliceToAny(productIds)...).Encode()
	})
}

type orderbookWrapper struct {
	Pricebook *Orderbook `json:"pricebook"`
}

type Orderbook struct {
	ProductId string    `json:"product_id"`
	Bids      []*Tick   `json:"bids"`
	Asks      []*Tick   `json:"asks"`
	Time      time.Time `json:"time"`
}

const ADVANCED_TRADE_ORDERBOOK_URL = "https://api.coinbase.com/api/v3/brokerage/product_book"

func (c *Client) GetOrderbook(ctx context.Context, productId string, limit int) (*Orderbook, error) {
	wrapper, err := http_helpers.GetJSONFn[*orderbookWrapper](ctx, c.cl, ADVANCED_TRADE_ORDERBOOK_URL, nil, func(r *http.Request) {
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).
			Add("product_id", productId).
			Add("limit", limit, limit != 0).
			Encode()
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

type Market struct {
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
	MidMarketPrice            float64     `json:"mid_market_price"`
	AliasTo                   []string    `json:"alias_to"`
	BaseDisplaySymbol         string      `json:"base_display_symbol"`
	QuoteDisplaySymbol        string      `json:"quote_display_symbol"`
	ViewOnly                  bool        `json:"view_only"`
	PriceIncrement            float64     `json:"price_increment,string"`
}

func (a *Market) UnmarshalJSON(bts []byte) error {
	if a == nil {
		a = &Market{}
	}
	var wrapped marketWrapper
	// ignore because error will occur due to empty string
	_ = json.Unmarshal(bts, &wrapped)
	*a = *unwrapMarket(&wrapped)
	return nil
}

// A bug in API returning empty strings instead of zero results in unmarshaler errors
// So instead we just marshal into this and then port fields over to result structed
// Against defining a custom float64 type to avoid requiring casting anywhere it is used
type marketWrapper struct {
	ProductId                 string      `json:"product_id"`
	Price                     string      `json:"price"`
	PricePercentageChange24H  string      `json:"price_percentage_change_24h"`
	Volume24H                 string      `json:"volume_24h"`
	VolumePercentageChange24H string      `json:"volume_percentage_change_24h"`
	BaseIncrement             string      `json:"base_increment"`
	QuoteIncrement            string      `json:"quote_increment"`
	QuoteMinSize              string      `json:"quote_min_size"`
	QuoteMaxSize              string      `json:"quote_max_size"`
	BaseMinSize               string      `json:"base_min_size"`
	BaseMaxSize               string      `json:"base_max_size"`
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
	PriceIncrement            string      `json:"price_increment"`
}

func unwrapMarket(wrapped *marketWrapper) *Market {
	unwrapped := &Market{
		ProductId:                 wrapped.ProductId,
		BaseName:                  wrapped.BaseName,
		Price:                     fastfloat.ParseBestEffort(wrapped.Price),
		PricePercentageChange24H:  fastfloat.ParseBestEffort(wrapped.PricePercentageChange24H),
		Volume24H:                 fastfloat.ParseBestEffort(wrapped.Volume24H),
		VolumePercentageChange24H: fastfloat.ParseBestEffort(wrapped.VolumePercentageChange24H),
		BaseIncrement:             fastfloat.ParseBestEffort(wrapped.BaseIncrement),
		QuoteIncrement:            fastfloat.ParseBestEffort(wrapped.QuoteIncrement),
		QuoteMinSize:              fastfloat.ParseBestEffort(wrapped.QuoteMinSize),
		QuoteMaxSize:              fastfloat.ParseBestEffort(wrapped.QuoteMaxSize),
		BaseMinSize:               fastfloat.ParseBestEffort(wrapped.BaseMinSize),
		BaseMaxSize:               fastfloat.ParseBestEffort(wrapped.BaseMaxSize),
		QuoteName:                 wrapped.QuoteName,
		Watched:                   wrapped.Watched,
		IsDisabled:                wrapped.IsDisabled,
		New:                       wrapped.New,
		Status:                    wrapped.Status,
		CancelOnly:                wrapped.CancelOnly,
		LimitOnly:                 wrapped.LimitOnly,
		PostOnly:                  wrapped.PostOnly,
		TradingDisabled:           wrapped.TradingDisabled,
		AuctionMode:               wrapped.AuctionMode,
		ProductType:               wrapped.ProductType,
		QuoteCurrencyId:           wrapped.QuoteCurrencyId,
		BaseCurrencyId:            wrapped.BaseCurrencyId,
		FcmTradingSessionDetails:  wrapped.FcmTradingSessionDetails,
		MidMarketPrice:            fastfloat.ParseBestEffort(wrapped.MidMarketPrice),
		AliasTo:                   wrapped.AliasTo,
		BaseDisplaySymbol:         wrapped.BaseDisplaySymbol,
		QuoteDisplaySymbol:        wrapped.QuoteDisplaySymbol,
		ViewOnly:                  wrapped.ViewOnly,
		PriceIncrement:            fastfloat.ParseBestEffort(wrapped.PriceIncrement),
	}
	return unwrapped
}

type Markets struct {
	Products []*Market `json:"products"`
}

const ADVANCED_TRADE_GET_MARKETS_URL = "https://api.coinbase.com/api/v3/brokerage/products"

func (c *Client) GetMarkets(ctx context.Context, params *GetMarketParams) (*Markets, error) {
	return http_helpers.GetJSONFn[*Markets](ctx, c.cl, ADVANCED_TRADE_GET_MARKETS_URL, nil, func(r *http.Request) {
		if params == nil {
			return
		}
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).
			AddCond("limit", params.Limit, params.Limit != 0).
			AddCond("offset", params.Offset, params.Offset != 0).
			AddCond("product_type", string(params.ProductType), params.ProductType != "").
			Add("product_ids", lambda.SliceToAny(params.ProductIds)...).
			AddCond("contract_expiry_type", string(params.ContractExpiryType), params.ContractExpiryType != "").
			Encode()
	})
}

const ADVANCED_TRADE_GET_MARKET_URL = "https://api.coinbase.com/api/v3/brokerage/products/%s"

func (c *Client) GetMarket(ctx context.Context, marketId string) (*Market, error) {
	return http_helpers.GetJSON[*Market](ctx, c.cl, fmt.Sprintf(ADVANCED_TRADE_GET_MARKET_URL, marketId), nil)
}

type Candle struct {
	Start  int64   `json:"start,string"`
	High   float64 `json:"high,string"`
	Low    float64 `json:"low,string"`
	Open   float64 `json:"open,string"`
	Close  float64 `json:"close,string"`
	Volume float64 `json:"volume,string"`
}

type Candles struct {
	Candles []*Candle `json:"candles"`
}

const ADVANCED_TRADE_CANDLES_URL = "https://api.coinbase.com/api/v3/brokerage/products/%s/candles"

type CandleGranularity string

var (
	CANDLE_GRANULARITY_1_MINUTE  CandleGranularity = "ONE_MINUTE"
	CANDLE_GRANULARITY_5_MINUTE  CandleGranularity = "FIVE_MINUTE"
	CANDLE_GRANULARITY_15_MINUTE CandleGranularity = "FIFTEEN_MINUTE"
	CANDLE_GRANULARITY_30_MINUTE CandleGranularity = "THIRTY_MINUTE"
	CANDLE_GRANULARITY_1_HOUR    CandleGranularity = "ONE_HOUR"
	CANDLE_GRANULARITY_2_HOUR    CandleGranularity = "TWO_HOUR"
	CANDLE_GRANULARITY_6_HOUR    CandleGranularity = "SIX_HOUR"
	CANDLE_GRANULARITY_1_DAY     CandleGranularity = "ONE_DAY"
)

func (c *Client) GetCandles(ctx context.Context, marketId string, granularity CandleGranularity, start, end int64) (*Candles, error) {
	return http_helpers.GetJSONFn[*Candles](ctx, c.cl, fmt.Sprintf(ADVANCED_TRADE_CANDLES_URL, marketId), nil, func(r *http.Request) {
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).
			Add("granularity", granularity).
			AddCond("start", start, start != 0 && end != 0).
			AddCond("end", end, start != 0 && end != 0).Encode()
	})
}

type MarketTrade struct {
	TradeId   int64     `json:"trade_id,string"`
	ProductId string    `json:"product_id"`
	Price     float64   `json:"price,string"`
	Size      float64   `json:"size,string"`
	Time      time.Time `json:"time"`
	Side      string    `json:"side"`
	Bid       string    `json:"bid"`
	Ask       string    `json:"ask"`
}

type MarketTrades struct {
	Trades []*MarketTrade
}

const ADVANCED_TRADE_MARKET_TRADES_URL = "https://api.coinbase.com/api/v3/brokerage/products/%s/ticker"

func (c *Client) GetMarketTrades(ctx context.Context, marketId string, limit int) (*MarketTrades, error) {
	return http_helpers.GetJSONFn[*MarketTrades](ctx, c.cl, fmt.Sprintf(ADVANCED_TRADE_MARKET_TRADES_URL, marketId), nil, func(r *http.Request) {
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).
			AddCond("limit", limit, limit != 0).
			Encode()
	})
}
