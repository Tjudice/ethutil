package coinbase

import (
	"net/http"
)

type Client struct {
	cl   *http.Client
	auth *AccountAuth
}

func NewClient(auth *AccountAuth) *Client {
	return &Client{
		cl:   http.DefaultClient,
		auth: auth,
	}
}
