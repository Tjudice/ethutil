package coinbase

import "net/http"

type AdvancedTradeClient struct {
	cl   *http.Client
	auth *AdvancedTradeAuth
}

func NewAdvancedTradeClient(auth *AdvancedTradeAuth) *AdvancedTradeClient {
	return &AdvancedTradeClient{
		cl:   http.DefaultClient,
		auth: auth,
	}
}
