package advanced_trade

import (
	"net/http"

	"github.com/tjudice/ethutil/exchange/coinbase/auth"
)

type Client struct {
	cl   *http.Client
	auth auth.Authenticator
}

const ADVANCED_TRADE_PATH_PREFIX = "https://api.coinbase.com"

func NewClient(a auth.Authenticator) *Client {
	return &Client{
		cl: &http.Client{
			Transport: &auth.Transport{
				Auth:         a,
				PrefixString: ADVANCED_TRADE_PATH_PREFIX,
			},
		},
		auth: a,
	}
}
