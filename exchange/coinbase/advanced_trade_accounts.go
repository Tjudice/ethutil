package coinbase

import (
	"context"
	"time"

	"github.com/tjudice/util/go/clients/jsonhttp"
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
	return jsonhttp.Get[*AccountsWrapper](ctx, c.cl, ADVANCED_TRADE_ACCOUNTS_URL, nil)
}
