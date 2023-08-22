package coinbase

import (
	"net/http"
)

type Client struct {
	cl   *http.Client
	auth Authenticator
}

func NewClient(auth Authenticator) *Client {
	return &Client{
		cl:   http.DefaultClient,
		auth: auth,
	}
}
