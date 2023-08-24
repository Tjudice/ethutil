package coinbase

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/tjudice/util/go/network/http_helpers"
)

type OrderStatus string

var (
	OrderStatusOpen      OrderStatus = "OPEN"
	OrderStatusFilled    OrderStatus = "FILLED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type OrderType string

var (
	OrderTypeMarket    OrderType = "MARKET"
	OrderTypeLimit     OrderType = "LIMIT"
	OrderTypeStop      OrderType = "STOP"
	OrderTypeStopLimit OrderType = "STOP_LIMIT"
)

type OrderSide string

var (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type ProductType string

var (
	ProductTypeSpot   ProductType = "SPOT"
	ProductTypeFuture ProductType = "FUTURE"
)

type OrderPlacementSource string

var (
	OrderPlacementSourceRetailSimple   OrderPlacementSource = "RETAIL_SIMPLE"
	OrderPlacementSourceRetailAdvanced OrderPlacementSource = "RETAIL_ADVANCED"
)

type ContractExpiryType string

var (
	ContractExpiryTypeExpiring ContractExpiryType = "EXPIRING"
)

type OrderParams struct {
	ProductId            string
	OrderStatus          []string
	Limit                int
	StartDate            time.Time
	EndDate              time.Time
	UserNativeCurrency   string
	OrderType            OrderType
	OrderSide            OrderSide
	Cursor               string
	ProductType          ProductType
	OrderPlacementSource OrderPlacementSource
	ContractExpiryType   ContractExpiryType
}

type Orders struct{}

type AdvancedTradeOrder struct{}

const ADVANCED_TRADE_ORDERS_URL = "https://api.coinbase.com/api/v3/brokerage/orders/historical/batch"

func (c *AdvancedTradeClient) GetOrders(ctx context.Context, params *OrderParams) (json.RawMessage, error) {
	return http_helpers.GetJSONFn[json.RawMessage](ctx, c.cl, ADVANCED_TRADE_ORDERS_URL, nil, func(r *http.Request) {
		if params == nil {
			return
		}
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).Add("product_id", params.ProductId).
			AddCond("order_status", params.OrderStatus, len(params.OrderStatus) != 0).
			AddCond("limit", params.Limit, params.Limit != 0).
			AddCond("start_date", params.StartDate, params.StartDate != time.Time{}).
			AddCond("end_date", params.EndDate, params.EndDate != time.Time{}).
			AddCond("user_native_currency", params.UserNativeCurrency, params.UserNativeCurrency != "").
			AddCond("order_type", params.OrderType, params.OrderType != "").
			AddCond("order_side", params.OrderSide, params.OrderSide != "").
			AddCond("cursor", params.Cursor, params.Cursor != "").
			AddCond("product_type", params.ProductType, params.ProductType != "").
			AddCond("order_placement_source", params.OrderPlacementSource, params.OrderPlacementSource != "").
			AddCond("contract_expiry_type", params.ContractExpiryType, params.ContractExpiryType != "").
			Encode()
	})
}

const ADVANCED_TRADE_ORDER_URL = "https://api.coinbase.com/api/v3/brokerage/orders/historical/%s"

func (c *AdvancedTradeClient) GetOrder(ctx context.Context, orderId string) (*AdvancedTradeOrder, error) {
	panic("not implemented")
}

type FillParams struct {
	OrderId                string
	ProductId              string
	StartSequenceTimestamp time.Time
	EndSequenceTimestamp   time.Time
	Limit                  int
	Cursor                 string
}

type Fills struct{}

const ADVANCED_TRADE_FILLS_URL = "https://api.coinbase.com/api/v3/brokerage/orders/historical/fills"

func (c *AdvancedTradeClient) GetFills(ctx context.Context, params *FillParams) (*Fills, error) {
	panic("not implemented")
}

type TransactionsSummaryParams struct {
	StartDate          time.Time
	EndDate            time.Time
	UserNativeCurrency string
	ProductType        ProductType
	ContractExpiryType ContractExpiryType
}

type TransactionsSummary struct{}

const ADVANCED_TRADE_TRANSACTIONS_SUMMARY_URL = "https://api.coinbase.com/api/v3/brokerage/transaction_summary"

func (c *AdvancedTradeClient) GetTransactionsSummary(ctx context.Context, params *TransactionsSummaryParams) (*TransactionsSummary, error) {
	panic("not implemented")
}
