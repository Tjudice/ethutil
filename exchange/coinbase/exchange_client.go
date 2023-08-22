package coinbase

import (
	"net/http"
)

type ExchangeClient struct {
	cl   *http.Client
	auth Authenticator
}

func NewExchangeClient(auth Authenticator) *ExchangeClient {
	return &ExchangeClient{
		cl: &http.Client{
			Transport: &Transport{
				Auth:         auth,
				PrefixString: "https://api.exchange.coinbase.com",
			},
		},
		auth: auth,
	}
}
