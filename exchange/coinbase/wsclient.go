package coinbase

import (
	"context"

	"github.com/gorilla/websocket"
)

type Conn struct {
	ch   chan WebsocketMessage
	conn *websocket.Conn
	auth *AccountAuth
}

const COINBASE_EXCHANGE_WSS_URL = "wss://ws-direct.exchange.coinbase.com"

type SubscribeMsg struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []any    `json:"channels"`
}

type SignedSubscribeMsg struct {
	*SubscribeMsg
	*SignedMessage
}

type WebsocketMessage interface{}

func (c *Client) Subscribe(ctx context.Context, products []string, channels []any) (*Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(COINBASE_EXCHANGE_WSS_URL, nil)
	if err != nil {
		return nil, err
	}
	msg := &SubscribeMsg{
		Type:       "subscribe",
		ProductIds: products,
		Channels:   channels,
	}
	signed, err := SignWebsocket(c.auth)
	if err != nil {
		return nil, err
	}
	err = conn.WriteJSON(&SignedSubscribeMsg{
		SubscribeMsg:  msg,
		SignedMessage: signed,
	})
	if err != nil {
		return nil, err
	}
	wsConn := &Conn{
		conn: conn,
		auth: c.auth,
		ch:   make(chan WebsocketMessage),
	}
	go wsConn.Listen()
	if err != nil {
		return nil, err
	}
	return wsConn, nil
}

func (c *Conn) Listen() error {
	defer c.conn.Close()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}
		parsed, err := parseMessage(msg)
		if err != nil {
			return err
		}
		c.ch <- parsed
	}
}

func (c *Conn) C() <-chan WebsocketMessage {
	return c.ch
}

func parseMessage(bts []byte) (WebsocketMessage, error) {
	return nil, nil
}
