package advanced_trade

const (
	ADVANCED_TRADE_WEBSOCKET_URL = "wss://advanced-trade-ws.coinbase.com"
)

type Websocket struct{}

func (c *Client) Subscribe() (*Websocket, error) {
	return nil, nil
}
