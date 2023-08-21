package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tjudice/util/go/clients/jsonhttp"
)

type Market struct {
	// Identifaction / Market Name
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	// Base Currency Info
	BaseCurrency  string  `json:"base_currency"`
	BaseIncrement float64 `json:"base_increment,string"`
	// Quote Currency Info
	QuoteCurrency  string  `json:"quote_currency"`
	QuoteIncrement float64 `json:"quote_increment,string"`
	// Market Parameters
	MinMarketFunds         float64 `json:"min_market_funds,string"`
	TradingDisabled        bool    `json:"trading_disabled"`
	MarginEnabled          bool    `json:"margin_enabled"`
	PostOnly               bool    `json:"post_only"`
	LimitOnly              bool    `json:"limit_only"`
	CancelOnly             bool    `json:"cancel_only"`
	MaxSlippage            float64 `json:"max_slippage_percentage,string"`
	FxStablecoin           bool    `json:"fx_stablecoin"`
	Status                 string  `json:"status"`
	StatusMessage          string  `json:"status_message"`
	AuctionMode            bool    `json:"auction_mode"`
	HighBidLimitPercentage string  `json:"high_bid_limit_percentage"`
}

const PRODUCTS_URL = "https://api.exchange.coinbase.com/products/"

func (c *Client) GetMarkets(ctx context.Context) ([]*Market, error) {
	return jsonhttp.Get[[]*Market](ctx, c.cl, PRODUCTS_URL, nil)
}

func (c *Client) GetMarket(ctx context.Context, marketId string) (*Market, error) {
	return jsonhttp.Get[*Market](ctx, c.cl, PRODUCTS_URL+marketId, nil)
}

const BOOK_URL = "https://api.exchange.coinbase.com/products/%s/book?level=%d"

type Order struct {
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
	NumOrders int     `json:"num_orders"`
}

type OrderbookWrapper struct {
	Sequence    int64     `json:"sequence"`
	Time        time.Time `json:"time"`
	AuctionMode bool      `json:"auction_mode"`
	Auction     *string   `json:"auction"`
}

type Orderbook struct {
	OrderbookWrapper
	Bids OrderMarshalling `json:"bids"`
	Asks OrderMarshalling `json:"asks"`
}

type OrderMarshalling []*Order

func (o *OrderMarshalling) UnmarshalJSON(bts []byte) error {
	var orderArray []json.RawMessage
	err := json.Unmarshal(bts, &orderArray)
	if err != nil {
		return err
	}
	var next [3]json.RawMessage
	orders := make([]*Order, 0, len(orderArray))
	for _, currOrder := range orderArray {
		err := json.Unmarshal(currOrder, &next)
		if err != nil {
			return err
		}
		nextOrder := &Order{}
		if err := unmarshalFloatString(next[0], &nextOrder.Price); err != nil {
			return err
		}
		if err := unmarshalFloatString(next[1], &nextOrder.Amount); err != nil {
			return err
		}
		if err := json.Unmarshal(next[2], &nextOrder.NumOrders); err != nil {
			return err
		}
		orders = append(orders, nextOrder)
	}
	*o = orders
	return nil
}

func (c *Client) GetMarketBookLevel1(ctx context.Context, marketId string) (*Orderbook, error) {
	return jsonhttp.Get[*Orderbook](ctx, c.cl, makeBookURL(marketId, 1), nil)
}

func (c *Client) GetMarketBookLevel2(ctx context.Context, marketId string) (*Orderbook, error) {
	return jsonhttp.Get[*Orderbook](ctx, c.cl, makeBookURL(marketId, 2), nil)
}

type OrderbookLevel3 struct {
	OrderbookWrapper
	Bids Order3Marshalling `json:"bids"`
	Asks Order3Marshalling `json:"asks"`
}

type Order3Marshalling []*Order3

type Order3 struct {
	Amount  float64 `json:"amount"`
	Price   float64 `json:"price"`
	OrderID string  `json:"order_id"`
}

func (o *Order3Marshalling) UnmarshalJSON(bts []byte) error {
	var orderArray []json.RawMessage
	err := json.Unmarshal(bts, &orderArray)
	if err != nil {
		return err
	}
	var next [3]json.RawMessage
	orders := make([]*Order3, 0, len(orderArray))
	for _, currOrder := range orderArray {
		err := json.Unmarshal(currOrder, &next)
		if err != nil {
			return err
		}
		nextOrder := &Order3{}
		if err := unmarshalFloatString(next[0], &nextOrder.Price); err != nil {
			return err
		}
		if err := unmarshalFloatString(next[1], &nextOrder.Amount); err != nil {
			return err
		}
		if err := json.Unmarshal(next[2], &nextOrder.OrderID); err != nil {
			return err
		}
		orders = append(orders, nextOrder)
	}
	*o = orders
	return nil
}

func (c *Client) GetMarketBookLevel3(ctx context.Context, marketId string) (*OrderbookLevel3, error) {
	return jsonhttp.Get[*OrderbookLevel3](ctx, c.cl, makeBookURL(marketId, 3), nil)
}

func makeBookURL(marketId string, level int) string {
	return fmt.Sprintf(BOOK_URL, marketId, level)
}

type Candles []*Candle

func (c *Candles) UnmarshalJSON(bts []byte) error {
	var candleArray []json.RawMessage
	err := json.Unmarshal(bts, &candleArray)
	if err != nil {
		return err
	}
	var buffer [6]json.RawMessage
	candles := make([]*Candle, 0, len(candleArray))
	for _, currCandle := range candleArray {
		err := json.Unmarshal(currCandle, &buffer)
		if err != nil {
			return err
		}
		nextCandle := &Candle{}
		if err := json.Unmarshal(buffer[0], &nextCandle.Ts); err != nil {
			return err
		}
		if err := json.Unmarshal(buffer[1], &nextCandle.Low); err != nil {
			return err
		}
		if err := json.Unmarshal(buffer[2], &nextCandle.High); err != nil {
			return err
		}
		if err := json.Unmarshal(buffer[3], &nextCandle.Open); err != nil {
			return err
		}
		if err := json.Unmarshal(buffer[4], &nextCandle.Close); err != nil {
			return err
		}
		if err := json.Unmarshal(buffer[5], &nextCandle.Volume); err != nil {
			return err
		}
		candles = append(candles, nextCandle)
	}
	*c = candles
	return nil
}

type Candle struct {
	Ts     int64   `json:"ts"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

const MARKET_CANDLES_URL = "https://api.exchange.coinbase.com/products/%s/candles?granularity=%d"

func (c *Client) GetMarketCandles(ctx context.Context, marketId string, granularity, start, end int) (*Candles, error) {
	return jsonhttp.Get[*Candles](ctx, c.cl, makeCandleURL(marketId, granularity, start, end), nil)
}

func makeCandleURL(marketId string, granularity, start, end int) string {
	url := fmt.Sprintf(MARKET_CANDLES_URL, marketId, granularity)
	if start != 0 && end != 0 {
		url = url + "&start=" + strconv.FormatInt(int64(start), 10) + "&end=" + strconv.FormatInt(int64(end), 10)
	}
	return url
}

type Stats struct {
	Open      float64 `json:"open,string"`
	High      float64 `json:"high,string"`
	Low       float64 `json:"low,string"`
	Last      float64 `json:"last,string"`
	Volume    float64 `json:"volume,string"`
	Volume30D float64 `json:"volume_30day,string"`
}

const MARKET_STATS_URL = "https://api.exchange.coinbase.com/products/%s/stats"

func (c *Client) GetMarketStats(ctx context.Context, marketId string) (*Stats, error) {
	return jsonhttp.Get[*Stats](ctx, c.cl, fmt.Sprintf(MARKET_STATS_URL, marketId), nil)
}

type Ticker struct {
	TradeId int64     `json:"trade_id"`
	Ask     float64   `json:"ask,string"`
	Bid     float64   `json:"bid,string"`
	Volume  float64   `json:"volume,string"`
	Price   float64   `json:"price,string"`
	Size    float64   `json:"size,string"`
	Time    time.Time `json:"time"`
}

const MARKET_TICKER_URL = "https://api.exchange.coinbase.com/products/%s/ticker"

func (c *Client) GetMarketTicker(ctx context.Context, marketId string) (*Ticker, error) {
	return jsonhttp.Get[*Ticker](ctx, c.cl, fmt.Sprintf(MARKET_TICKER_URL, marketId), nil)
}
