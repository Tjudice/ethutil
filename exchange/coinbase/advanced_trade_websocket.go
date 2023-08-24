package coinbase

const (
	ADVANCED_TRADE_WEBSOCKET_URL = "wss://advanced-trade-ws.coinbase.com"
)

type AdvancedTradeWebsocket struct{}

func (c *AdvancedTradeClient) Subscribe() (*AdvancedTradeWebsocket, error) {
	return nil, nil
}
