package coinbase

import (
	"context"

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

const PRODUCTS_URL = "https://api.exchange.coinbase.com/products"

func (c *Client) GetMarkets(ctx context.Context) ([]*Market, error) {
	res, err := jsonhttp.Get[[]*Market](ctx, c.cl, PRODUCTS_URL, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}
