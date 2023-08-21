package coinbase

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

type Conn struct {
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

func Subscribe(ctx context.Context, a *AccountAuth, products []string, channels []any) (*Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(COINBASE_EXCHANGE_WSS_URL, nil)
	if err != nil {
		return nil, err
	}
	msg := &SubscribeMsg{
		Type:       "subscribe",
		ProductIds: products,
		Channels:   channels,
	}
	signed, err := SignWebsocket(a)
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
	c := &Conn{
		conn: conn,
		auth: a,
	}
	go c.Listen()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Conn) Listen() error {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		log.Println(string(msg))
	}
}
