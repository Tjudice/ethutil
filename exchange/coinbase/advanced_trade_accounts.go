package coinbase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tjudice/util/go/network/http_helpers"
)

type AccountsWrapper struct {
	Accounts []*Account `json:"accounts"`
	HasNext  bool       `json:"has_next"`
	Cursor   string     `json:"cursor"`
	Size     int        `json:"size"`
}

type Account struct {
	UUID             string    `json:"uuid"`
	Name             string    `json:"name"`
	Currency         string    `json:"currency"`
	AvailableBalance Balance   `json:"available_balance"`
	Default          bool      `json:"default"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        *string   `json:"deleted_at"`
	Type             string    `json:"type"`
	Ready            bool      `json:"ready"`
	Hold             Balance   `json:"hold"`
}

type Balance struct {
	Value    float64 `json:"value,string"`
	Currency string  `json:"currency"`
}

var ADVANCED_TRADE_ACCOUNTS_URL = "https://api.coinbase.com/api/v3/brokerage/accounts"

func (c *AdvancedTradeClient) GetAccounts(ctx context.Context, limit int, cursor string) (*AccountsWrapper, error) {
	return http_helpers.GetJSONFn[*AccountsWrapper](ctx, c.cl, ADVANCED_TRADE_ACCOUNTS_URL, nil, func(r *http.Request) {
		r.URL.RawQuery, _ = http_helpers.NewURLEncoder(r.URL.Query()).
			AddIfNotDefault("limit", limit, 0).
			AddIfNotDefault("cursor", cursor, "").Encode()
	})
}

type AccountWrapper struct {
	Account *Account `json:"account"`
}

var ADVANCED_TRADE_ACCOUNT_URL = "https://api.coinbase.com/api/v3/brokerage/accounts/%s"

func (c *AdvancedTradeClient) GetAccount(ctx context.Context, accountUUID string) (*Account, error) {
	acc, err := http_helpers.GetJSON[*AccountWrapper](ctx, c.cl, fmt.Sprintf(ADVANCED_TRADE_ACCOUNT_URL, accountUUID), nil)
	if err != nil {
		return nil, err
	}
	return acc.Account, nil
}
