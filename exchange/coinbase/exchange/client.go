package exchange

import (
	"net/http"

	"github.com/tjudice/ethutil/exchange/coinbase/auth"
)

type Client struct {
	cl   *http.Client
	auth auth.Authenticator
}

func NewClient(a auth.Authenticator) *Client {
	return &Client{
		cl: &http.Client{
			Transport: &auth.Transport{
				Auth:         a,
				PrefixString: "https://api.exchange.coinbase.com",
			},
		},
		auth: a,
	}
}
