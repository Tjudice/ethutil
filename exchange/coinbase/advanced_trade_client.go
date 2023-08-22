package coinbase

import "net/http"

type AdvancedTradeClient struct {
	cl   *http.Client
	auth Authenticator
}

const ADVANCED_TRADE_PATH_PREFIX = "https://api.coinbase.com"

func NewAdvancedTradeClient(auth Authenticator) *AdvancedTradeClient {
	return &AdvancedTradeClient{
		cl: &http.Client{
			Transport: &Transport{
				Auth:         auth,
				PrefixString: ADVANCED_TRADE_PATH_PREFIX,
			},
		},
		auth: auth,
	}
}
