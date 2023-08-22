package coinbase

import (
	"context"
	"time"
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

func (c *AdvancedTradeClient) GetOrders(ctx context.Context, params *OrderParams) (*Orders, error) {
	panic("not implemented")
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
