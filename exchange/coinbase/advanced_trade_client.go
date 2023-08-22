package coinbase

import "net/http"

type AdvancedTradeClient struct {
	cl   *http.Client
	auth Authenticator
}

func NewAdvancedTradeClient(auth Authenticator) *AdvancedTradeClient {
	return &AdvancedTradeClient{
		cl:   http.DefaultClient,
		auth: auth,
	}
}
