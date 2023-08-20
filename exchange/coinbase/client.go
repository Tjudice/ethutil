package coinbase

import (
	"net/http"
)

type Client struct {
	cl *http.Client
}

func NewClient() *Client {
	return &Client{
		cl: http.DefaultClient,
	}
}
