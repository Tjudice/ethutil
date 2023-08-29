package advanced_trade

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tjudice/ethutil/exchange/coinbase/auth"
)

const (
	ADVANCED_TRADE_WEBSOCKET_URL = "wss://advanced-trade-ws.coinbase.com"
)

type Channel string

const (
	SubscriptionsChannel Channel = "subscriptions"
	HeartbeatsChannel    Channel = "heartbeats"
	CandlesChannel       Channel = "candles"
)

type SubscribeMsg struct {
	Type       string   `json:"type"`
	Channel    Channel  `json:"channel"`
	ProductIds []string `json:"product_ids"`
	Sig
}

type Conn struct {
	conn *websocket.Conn
	auth auth.Authenticator
	ch   chan WebsocketMessage
	done chan struct{}
}

type WebsocketMessage interface {
	Seq() int64
	Clone() WebsocketMessage
}

func (c *Client) Subscribe(ctx context.Context, bufferSize int, channel Channel, productIds []string) (*Conn, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, ADVANCED_TRADE_WEBSOCKET_URL, nil)
	if err != nil {
		return nil, err
	}
	msg := &SubscribeMsg{
		Type:       "subscribe",
		Channel:    channel,
		ProductIds: productIds,
	}
	signed, err := c.auth.SignWebsocketRequest([]string{string(channel)}, productIds)
	if err != nil {
		return nil, err
	}
	msg.Sig = Sig{
		Key:       signed.Key,
		Signature: signed.Sig,
		Timestamp: signed.Timestamp,
	}
	log.Println(msg)
	err = conn.WriteJSON(msg)
	if err != nil {
		return nil, err
	}
	wsConn := &Conn{
		conn: conn,
		auth: c.auth,
		ch:   make(chan WebsocketMessage, bufferSize),
		done: make(chan struct{}),
	}
	go wsConn.listen()
	if err != nil {
		return nil, err
	}
	return wsConn, nil
}

func (c *Conn) listen() {
	defer c.conn.Close()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		parsed, err := parseMessage(msg)
		if err != nil {
			log.Println(err.Error())
			return
		}
		c.ch <- parsed
	}
}

func (c *Conn) C() <-chan WebsocketMessage {
	return c.ch
}

var messageTypeChoice = map[Channel]WebsocketMessage{
	SubscriptionsChannel: &Subscriptions{},
	HeartbeatsChannel:    &Heartbeats{},
	CandlesChannel:       &CandlesFeed{},
}

type MessageType struct {
	Channel Channel `json:"channel"`
}

func parseMessage(bts []byte) (WebsocketMessage, error) {
	msgType := MessageType{}
	err := json.Unmarshal(bts, &msgType)
	if err != nil {
		return nil, err
	}
	schem, ok := messageTypeChoice[msgType.Channel]
	if !ok {
		log.Println(msgType.Channel)
		log.Println(string(bts))
		panic("a")
	}
	typedMessage := schem.Clone()
	err = json.Unmarshal(bts, &typedMessage)
	if err != nil {
		log.Println(msgType.Channel)
		log.Println(string(bts))
		return nil, err
	}
	return typedMessage, nil
}

type MessageDetails struct {
	ClientId    string    `json:"client_id"`
	Timestamp   time.Time `json:"timestamp"`
	SequenceNum int64     `json:"sequence_num"`
}

type Subscriptions struct {
	MessageDetails
	Events json.RawMessage `json:"events"`
}

func (s *Subscriptions) Seq() int64 {
	return s.SequenceNum
}

func (s *Subscriptions) Clone() WebsocketMessage {
	return new(Subscriptions)
}

type Heartbeats struct {
	MessageDetails
	Events []*HeartbeatEvent `json:"events"`
}

type HeartbeatEvent struct {
	CurrentTime      string `json:"current_time"`
	HeartbeatCounter int64  `json:"heartbeat_counter"`
}

func (s *Heartbeats) Seq() int64 {
	return s.SequenceNum
}

func (s *Heartbeats) Clone() WebsocketMessage {
	return new(Heartbeats)
}

type CandlesFeed struct {
	MessageDetails
	Events []*CandleEvent `json:"events"`
}

type CandleEvent struct {
	Type    string        `json:"type"`
	Candles []*CandleData `json:"candles"`
}

type CandleData struct {
	ProductId string  `json:"product_id"`
	Start     int64   `json:"start,string"`
	Open      float64 `json:"open,string"`
	High      float64 `json:"high,string"`
	Low       float64 `json:"low,string"`
	Close     float64 `json:"close,string"`
	Volume    float64 `json:"volume,string"`
}

func (s *CandlesFeed) Seq() int64 {
	return s.SequenceNum
}

func (s *CandlesFeed) Clone() WebsocketMessage {
	return new(CandlesFeed)
}
