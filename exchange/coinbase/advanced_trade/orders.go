package advanced_trade

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tjudice/util/go/network/http_helpers"
	"github.com/valyala/fastjson/fastfloat"
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

func encodeOrderParams(r *http.Request, params *OrderParams) {
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
}

type Orders struct {
	Orders   []*AdvancedTradeOrder `json:"orders"`
	Sequence int64                 `json:"sequence,string"`
	HasNext  bool                  `json:"has_next"`
	Cursor   string                `json:"cursor"`
}

type OrderConfig interface {
	OrderType() OrderConfigDescriptor
}

type OrderConfiguration struct {
	ConfigurationType OrderConfigDescriptor `json:"-"`
	QuoteSize         float64               `json:"quote_size,string,omitempty"`
	BaseSize          float64               `json:"base_size,string,omitempty"`
	PostOnly          bool                  `json:"post_only,omitempty"`
	EndTime           time.Time             `json:"end_time,omitempty"`
	StopPrice         float64               `json:"stop_price,string,omitempty"`
	StopDirection     string                `json:"stop_direction,omitempty"`
	LimitPrice        float64               `json:"limit_price,string,omitempty"`
}

type orderConfigWrapper struct {
	QuoteSize     string    `json:"quote_size,omitempty"`
	BaseSize      string    `json:"base_size,omitempty"`
	PostOnly      bool      `json:"post_only,omitempty"`
	EndTime       time.Time `json:"end_time,omitempty"`
	StopPrice     string    `json:"stop_price,omitempty"`
	StopDirection string    `json:"stop_direction,omitempty"`
	LimitPrice    string    `json:"limit_price,omitempty"`
}

type OrderConfigDescriptor string

var (
	OrderConfigMarketMarketIOC       OrderConfigDescriptor = "market_market_ioc"
	OrderConfigLimitLimitGTC         OrderConfigDescriptor = "limit_limit_gtc"
	OrderConfigLimitLimitGTD         OrderConfigDescriptor = "limit_limit_gtd"
	OrderConfigStopLimitStopLimitGTC OrderConfigDescriptor = "stop_limit_stop_limit_gtc"
	OrderConfigStopLimitStopLimitGTD OrderConfigDescriptor = "stop_limit_stop_limit_gtd"
)

type MarketMarketIOC struct {
	QuoteSize float64 `json:"quote_size,string"`
	BaseSize  float64 `json:"base_size,string"`
}

func (m *MarketMarketIOC) OrderType() OrderConfigDescriptor {
	return OrderConfigMarketMarketIOC
}

type LimitLimitGTC struct {
	QuoteSize float64 `json:"quote_size,string"`
	BaseSize  float64 `json:"base_size,string"`
	PostOnly  bool    `json:"post_only"`
}

func (l *LimitLimitGTC) OrderType() OrderConfigDescriptor {
	return OrderConfigLimitLimitGTC
}

type LimitLimitGTD struct {
	BaseSize   float64   `json:"base_size,string"`
	LimitPrice float64   `json:"limit_price,string"`
	EndTime    time.Time `json:"end_time"`
	PostOnly   bool      `json:"post_only"`
}

func (l *LimitLimitGTD) OrderType() OrderConfigDescriptor {
	return OrderConfigLimitLimitGTD
}

type StopLimitStopLimitGTC struct {
	BaseSize      float64 `json:"base_size,string"`
	LimitPrice    float64 `json:"limit_price,string"`
	StopPrice     float64 `json:"stop_price,string"`
	StopDirection string  `json:"stop_direction"`
}

func (l *StopLimitStopLimitGTC) OrderType() OrderConfigDescriptor {
	return OrderConfigStopLimitStopLimitGTC
}

type StopLimitStopLimitGTD struct {
	BaseSize      float64   `json:"base_size"`
	LimitPrice    float64   `json:"limit_price,string"`
	StopPrice     float64   `json:"stop_price,string"`
	EndTime       time.Time `json:"end_time"`
	StopDirection string    `json:"stop_direction"`
}

func (l *StopLimitStopLimitGTD) OrderType() OrderConfigDescriptor {
	return OrderConfigStopLimitStopLimitGTD
}

type SafeOrderFields struct {
	OrderId              string               `json:"order_id"`
	ProductId            string               `json:"product_id"`
	UserId               string               `json:"user_id"`
	Side                 OrderSide            `json:"side"`
	ClientOrderId        string               `json:"client_order_id"`
	Status               OrderStatus          `json:"status"`
	TimeInForce          string               `json:"time_in_force"`
	CreatedTime          time.Time            `json:"created_time"`
	PendingCancel        bool                 `json:"pending_cancel"`
	SizeInQuote          bool                 `json:"size_in_quote"`
	SizeInclusiveOffees  bool                 `json:"size_inclusive_of_fees"`
	TriggerStatus        string               `json:"trigger_status"`
	OrderType            OrderType            `json:"order_type"`
	RejectReason         string               `json:"reject_reason"`
	Settled              bool                 `json:"settled"`
	ProductType          ProductType          `json:"product_type"`
	RejectMessage        string               `json:"reject_message"`
	OrderPlacementSource OrderPlacementSource `json:"order_placement_source"`
	IsLiquidation        bool                 `json:"is_liquidiation"`
	LastFillTime         time.Time            `json:"last_fill_time"`
	EditHistory          []string             `json:"edit_history"`
}

type customOrderFields struct {
	CompletetionPercentage string                                        `json:"completion_percentage"`
	FilledSize             string                                        `json:"filled_size"`
	AverageFilledPrice     string                                        `json:"average_filled_price"`
	Fee                    string                                        `json:"fee"`
	NumberOfFills          string                                        `json:"number_of_fills"`
	FilledValue            string                                        `json:"filled_value"`
	TotalFees              string                                        `json:"total_fees"`
	TotalValueAfterFees    string                                        `json:"total_value_after_fees"`
	OutstandingHoldAmount  string                                        `json:"outstanding_hold_amount"`
	OrderConfiguration     map[OrderConfigDescriptor]*orderConfigWrapper `json:"order_configuration"`
}

type AdvancedTradeOrder struct {
	CompletionPercentage  float64             `json:"completion_percentage,string"`
	FilledSize            float64             `json:"filled_size,string"`
	AverageFilledPrice    float64             `json:"average_filled_price,string"`
	Fee                   float64             `json:"fee,string"`
	NumberOfFills         int64               `json:"number_of_fills,string"`
	FilledValue           float64             `json:"filled_value,string"`
	TotalFees             float64             `json:"total_fees"`
	TotalValueAfterFees   float64             `json:"total_value_after_fees,string"`
	OutstandingHoldAmount float64             `json:"outstanding_hold_amount,string"`
	OrderConfiguration    *OrderConfiguration `json:"order_configuration"`
	SafeOrderFields
}

type orderWrapper struct {
	SafeOrderFields
	customOrderFields
}

func (a *AdvancedTradeOrder) UnmarshalJSON(bts []byte) error {
	if a == nil {
		a = &AdvancedTradeOrder{}
	}
	var wrapped orderWrapper
	err := json.Unmarshal(bts, &wrapped)
	if err != nil {
		return err
	}
	// This ignores the case if there are multiple order configurations. Almost positive this is not possible
	var conf OrderConfiguration
	for k, v := range wrapped.OrderConfiguration {
		conf = OrderConfiguration{
			ConfigurationType: OrderConfigDescriptor(k),
			BaseSize:          fastfloat.ParseBestEffort(v.BaseSize),
			QuoteSize:         fastfloat.ParseBestEffort(v.QuoteSize),
			PostOnly:          v.PostOnly,
			EndTime:           v.EndTime,
			StopPrice:         fastfloat.ParseBestEffort(v.StopPrice),
			StopDirection:     v.StopDirection,
			LimitPrice:        fastfloat.ParseBestEffort(v.LimitPrice),
		}
		break
	}
	*a = AdvancedTradeOrder{
		SafeOrderFields:       wrapped.SafeOrderFields,
		CompletionPercentage:  fastfloat.ParseBestEffort(wrapped.CompletetionPercentage),
		FilledSize:            fastfloat.ParseBestEffort(wrapped.FilledSize),
		AverageFilledPrice:    fastfloat.ParseBestEffort(wrapped.AverageFilledPrice),
		Fee:                   fastfloat.ParseBestEffort(wrapped.Fee),
		NumberOfFills:         fastfloat.ParseInt64BestEffort(wrapped.NumberOfFills),
		FilledValue:           fastfloat.ParseBestEffort(wrapped.FilledValue),
		TotalFees:             fastfloat.ParseBestEffort(wrapped.TotalFees),
		TotalValueAfterFees:   fastfloat.ParseBestEffort(wrapped.TotalValueAfterFees),
		OutstandingHoldAmount: fastfloat.ParseBestEffort(wrapped.OutstandingHoldAmount),
		OrderConfiguration:    &conf,
	}
	return nil
}

const ADVANCED_TRADE_ORDERS_URL = "https://api.coinbase.com/api/v3/brokerage/orders/historical/batch"

func (c *Client) GetOrders(ctx context.Context, params *OrderParams) (*Orders, error) {
	return http_helpers.GetJSONFn[*Orders](ctx, c.cl, ADVANCED_TRADE_ORDERS_URL, nil, func(r *http.Request) {
		encodeOrderParams(r, params)
	})
}

type singleOrderWrapper struct {
	Order *AdvancedTradeOrder `json:"order"`
}

const ADVANCED_TRADE_ORDER_URL = "https://api.coinbase.com/api/v3/brokerage/orders/historical/%s"

func (c *Client) GetOrder(ctx context.Context, orderId string) (*AdvancedTradeOrder, error) {
	wrapped, err := http_helpers.GetJSON[*singleOrderWrapper](ctx, c.cl, fmt.Sprintf(ADVANCED_TRADE_ORDER_URL, orderId), nil)
	if err != nil {
		return nil, err
	}
	return wrapped.Order, nil
}

type FillParams struct {
	OrderId                string
	ProductId              string
	StartSequenceTimestamp time.Time
	EndSequenceTimestamp   time.Time
	Limit                  int
	Cursor                 string
}

type Fills struct {
	Fills  []*Fill `json:"fills"`
	Cursor string  `json:"cursor"`
}

type Fill struct {
	EntryId            string    `json:"entry_id"`
	TradeId            string    `json:"trade_id"`
	OrderId            string    `json:"order_id"`
	TradeTime          time.Time `json:"trade_time"`
	TradeType          string    `json:"trade_type"`
	Price              float64   `json:"price,string"`
	Size               float64   `json:"size,string"`
	Commission         float64   `json:"commission,string"`
	ProductId          string    `json:"product_id"`
	SequenceTimestamp  time.Time `json:"sequence_timestamp"`
	LiquidityIndicator string    `json:"liquidity_indicator"`
	SizeInQuote        bool      `json:"size_in_quote"`
	UserId             string    `json:"user_id"`
	Side               OrderSide `json:"side"`
}

const ADVANCED_TRADE_FILLS_URL = "https://api.coinbase.com/api/v3/brokerage/orders/historical/fills"

func (c *Client) GetFills(ctx context.Context, params *FillParams) (*Fills, error) {
	return http_helpers.GetJSONFn[*Fills](ctx, c.cl, ADVANCED_TRADE_FILLS_URL, nil, func(r *http.Request) {
		if params == nil {
			return
		}
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).
			Add("product_id", params.ProductId).
			AddCond("order_id", params.OrderId, params.OrderId != "").
			AddCond("limit", params.Limit, params.Limit != 0).
			AddCond("cursor", params.Cursor, params.Cursor != "").
			AddCond("product_id", params.ProductId, params.ProductId != "").
			AddCond("start_sequence_timestamp", params.StartSequenceTimestamp, params.StartSequenceTimestamp != time.Time{}).
			AddCond("end_sequence_timestamp", params.EndSequenceTimestamp, params.EndSequenceTimestamp != time.Time{}).
			Encode()
	})
}

type TransactionsSummaryParams struct {
	StartDate          time.Time
	EndDate            time.Time
	UserNativeCurrency string
	ProductType        ProductType
	ContractExpiryType ContractExpiryType
}

type TransactionsSummary struct {
	TotalVolume             float64              `json:"total_volume"`
	TotalFees               float64              `json:"total_fees"`
	FeeTier                 *FeeTier             `json:"fee_tier"`
	MarginRate              *MarginRate          `json:"margin_rate"`
	GoodsAndServicesTax     *GoodsAndServicesTax `json:"goods_and_services_tax"`
	AdvancedTradeOnlyVolume float64              `json:"advanced_trade_only_volume"`
	AdvancedTradeOnlyFees   float64              `json:"advanced_trade_only_fees"`
	CoinbaseProVolume       float64              `json:"coinbase_pro_volume"`
	CoinbaseProFees         float64              `json:"coinbase_pro_fees"`
}

type FeeTier struct {
	PricingTier  string  `json:"pricing_tier"`
	UsdFrom      float64 `json:"usd_from"`
	UsdTo        float64 `json:"usd_to"`
	TakerFeeRate float64 `json:"taker_fee_rate"`
	MakerFeeRate float64 `json:"maker_fee_rate"`
}

type feeTierWrapper struct {
	PricingTier  string `json:"pricing_tier"`
	UsdFrom      string `json:"usd_from"`
	UsdTo        string `json:"usd_to"`
	TakerFeeRate string `json:"taker_fee_rate"`
	MakerFeeRate string `json:"maker_fee_rate"`
}

func (f *FeeTier) UnmarshalJSON(bts []byte) error {
	var wrapper feeTierWrapper
	err := json.Unmarshal(bts, &wrapper)
	if err != nil {
		return err
	}
	*f = FeeTier{
		PricingTier:  wrapper.PricingTier,
		UsdFrom:      fastfloat.ParseBestEffort(strings.Replace(wrapper.UsdFrom, ",", "", -1)),
		UsdTo:        fastfloat.ParseBestEffort(strings.Replace(wrapper.UsdTo, ",", "", -1)),
		TakerFeeRate: fastfloat.ParseBestEffort(wrapper.TakerFeeRate),
		MakerFeeRate: fastfloat.ParseBestEffort(wrapper.MakerFeeRate),
	}
	return nil
}

type MarginRate struct {
	Value string `json:"value"`
}

type GoodsAndServicesTax struct {
	Rate string `json:"rate"`
	Type string `json:"type"`
}

const ADVANCED_TRADE_TRANSACTIONS_SUMMARY_URL = "https://api.coinbase.com/api/v3/brokerage/transaction_summary"

func (c *Client) GetTransactionsSummary(ctx context.Context, params *TransactionsSummaryParams) (*TransactionsSummary, error) {
	return http_helpers.GetJSONFn[*TransactionsSummary](ctx, c.cl, ADVANCED_TRADE_TRANSACTIONS_SUMMARY_URL, nil, func(r *http.Request) {
		if params == nil {
			return
		}
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).
			AddCond("start_date", params.StartDate, params.StartDate != time.Time{}).
			AddCond("end_date", params.EndDate, params.EndDate != time.Time{}).
			AddCond("user_native_currency", params.UserNativeCurrency, params.UserNativeCurrency != "").
			AddCond("product_type", params.ProductType, params.ProductType != "").
			AddCond("contract_expiry_type", params.ContractExpiryType, params.ContractExpiryType != "").
			Encode()
	})
}
