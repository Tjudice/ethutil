package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	MidMarketPrice            float64     `json:"mid_market_price"`
	AliasTo                   []string    `json:"alias_to"`
	BaseDisplaySymbol         string      `json:"base_display_symbol"`
	QuoteDisplaySymbol        string      `json:"quote_display_symbol"`
	ViewOnly                  bool        `json:"view_only"`
	PriceIncrement            float64     `json:"price_increment,string"`
}

func (a *AdvancedTradeMarket) UnmarshalJSON(bts []byte) error {
	if a == nil {
		a = &AdvancedTradeMarket{}
	}
	var wrapped advancedTradeMarketWrapper
	// ignore because error will occur due to empty string
	_ = json.Unmarshal(bts, &wrapped)
	*a = *unwrapMarket(&wrapped)
	return nil
}

// A bug in API returning empty strings instead of zero results in unmarshaler errors
// So instead we just marshal into this and then port fields over to result structed
// Against defining a custom float64 type to avoid requiring casting anywhere it is used
type advancedTradeMarketWrapper struct {
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

func unwrapMarket(wrapped *advancedTradeMarketWrapper) *AdvancedTradeMarket {
	unwrapped := &AdvancedTradeMarket{
		ProductId:                 wrapped.ProductId,
		BaseName:                  wrapped.BaseName,
		Price:                     fastfloat.ParseBestEffort(strings.Trim(wrapped.Price, `"`)),
		PricePercentageChange24H:  fastfloat.ParseBestEffort(strings.Trim(wrapped.PricePercentageChange24H, `"`)),
		Volume24H:                 fastfloat.ParseBestEffort(strings.Trim(wrapped.Volume24H, `"`)),
		VolumePercentageChange24H: fastfloat.ParseBestEffort(strings.Trim(wrapped.VolumePercentageChange24H, `"`)),
		BaseIncrement:             fastfloat.ParseBestEffort(strings.Trim(wrapped.BaseIncrement, `"`)),
		QuoteIncrement:            fastfloat.ParseBestEffort(strings.Trim(wrapped.QuoteIncrement, `"`)),
		QuoteMinSize:              fastfloat.ParseBestEffort(strings.Trim(wrapped.QuoteMinSize, `"`)),
		QuoteMaxSize:              fastfloat.ParseBestEffort(strings.Trim(wrapped.QuoteMaxSize, `"`)),
		BaseMinSize:               fastfloat.ParseBestEffort(strings.Trim(wrapped.BaseMinSize, `"`)),
		BaseMaxSize:               fastfloat.ParseBestEffort(strings.Trim(wrapped.BaseMaxSize, `"`)),
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
	unwrapped.Price = fastfloat.ParseBestEffort(wrapped.Price)
	return unwrapped
}

type AdvancedTradeMarkets struct {
	Products []*AdvancedTradeMarket `json:"products"`
}

const ADVANCED_TRADE_GET_MARKETS_URL = "https://api.coinbase.com/api/v3/brokerage/products"

func (c *AdvancedTradeClient) GetMarkets(ctx context.Context, params *GetMarketParams) (*AdvancedTradeMarkets, error) {
	return http_helpers.GetJSONFn[*AdvancedTradeMarkets](ctx, c.cl, ADVANCED_TRADE_GET_MARKETS_URL, nil, func(r *http.Request) {
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
	return http_helpers.GetJSON[*AdvancedTradeMarket](ctx, c.cl, fmt.Sprintf(ADVANCED_TRADE_GET_MARKET_URL, marketId), nil)
}

type AdvancedTradeCandle struct {
	Start  int64   `json:"start,string"`
	High   float64 `json:"high,string"`
	Low    float64 `json:"low,string"`
	Open   float64 `json:"open,string"`
	Close  float64 `json:"close,string"`
	Volume float64 `json:"volume,string"`
}

type AdvancedTradeCandles struct {
	Candles []*AdvancedTradeCandle `json:"candles"`
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

func (c *AdvancedTradeClient) GetCandles(ctx context.Context, marketId string, granularity CandleGranularity, start, end int64) (*AdvancedTradeCandles, error) {
	return http_helpers.GetJSONFn[*AdvancedTradeCandles](ctx, c.cl, fmt.Sprintf(ADVANCED_TRADE_CANDLES_URL, marketId), nil, func(r *http.Request) {
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).Add("granularity", granularity).
			AddIfNotDefault("start", start, 0).AddIfNotDefault("end", end, 0).Encode()
	})
}

type AdvancedTradeMarketTrade struct{}

type AdvancedTradeMarketTrades struct{}

const ADVANCED_TRADE_MARKET_TRADES_URL = "https://api.coinbase.com/api/v3/brokerage/products/%s/ticker"

func (c *AdvancedTradeClient) GetMarketTrades(ctx context.Context, marketId string, limit int) (*AdvancedTradeMarketTrades, error) {
	panic("not implemented")
}
