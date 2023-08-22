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

func (c *AdvancedTradeClient) GetOrders(ctx context.Context, params *OrderParams) ([]*Order, error) {
	return nil, nil
}
