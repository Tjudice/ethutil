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
		cl:   http.DefaultClient,
		auth: auth,
	}
}
