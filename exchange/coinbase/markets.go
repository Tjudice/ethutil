package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tjudice/util/go/clients/jsonhttp"
)

type Market struct {
	// Identifaction / Market Name
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	// Base Currency Info
	BaseCurrency  string `json:"base_currency"`
	BaseIncrement string `json:"base_increment"`
	// Quote Currency Info
	QuoteCurrency  string `json:"quote_currency"`
	QuoteIncrement string `json:"quote_increment"`
	// Market Parameters
	MinMarketFunds         string `json:"min_market_funds"`
	TradingDisabled        bool   `json:"trading_disabled"`
	MarginEnabled          bool   `json:"margin_enabled"`
	PostOnly               bool   `json:"post_only"`
	LimitOnly              bool   `json:"limit_only"`
	CancelOnly             bool   `json:"cancel_only"`
	MaxSlippage            string `json:"max_slippage_percentage"`
	FxStablecoin           bool   `json:"fx_stablecoin"`
	Status                 string `json:"status"`
	StatusMessage          string `json:"status_message"`
	AuctionMode            bool   `json:"auction_mode"`
	HighBidLimitPercentage string `json:"high_bid_limit_percentage"`
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

type Orderbook struct {
	Sequence    int64            `json:"sequence"`
	Time        time.Time        `json:"time"`
	AuctionMode bool             `json:"auction_mode"`
	Auction     *string          `json:"auction"`
	Bids        OrderMarshalling `json:"bids"`
	Asks        OrderMarshalling `json:"asks"`
}

type OrderMarshalling []*Order

var _ json.Unmarshaler = &OrderMarshalling{}

func (o *OrderMarshalling) UnmarshalJSON(bts []byte) error {
	var orderArray []json.RawMessage
	err := json.Unmarshal(bts, &orderArray)
	if err != nil {
		return err
	}
	orders := make([]*Order, 0, len(orderArray))
	for _, currOrder := range orderArray {
		nextOrder := &Order{}
		orderVars := []interface{}{&FloatStringWrapper{&nextOrder.Price}, &FloatStringWrapper{&nextOrder.Amount}, &nextOrder.NumOrders}
		err := json.Unmarshal(currOrder, &orderVars)
		if err != nil {
			return err
		}
		orders = append(orders, nextOrder)
	}
	*o = orders
	return nil
}

func (c *Client) GetMarketBook(ctx context.Context, marketId string, level int) (*Orderbook, error) {
	return jsonhttp.Get[*Orderbook](ctx, c.cl, makeBookURL(marketId, level), nil)
}

func makeBookURL(marketId string, level int) string {
	return fmt.Sprintf(BOOK_URL, marketId, level)
}
