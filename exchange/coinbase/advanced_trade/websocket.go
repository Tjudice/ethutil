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
	MarketTradesChannel  Channel = "market_trades"
	StatusChannel        Channel = "status"
	TickerChannel        Channel = "ticker"
	TickerBatchChannel   Channel = "ticker_batch"
	Level2Channel        Channel = "level2"
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
	CandlesChannel:       &CandlesFeedItem{},
	MarketTradesChannel:  &MarketTradesFeedItem{},
	StatusChannel:        &StatusFeedItem{},
	TickerChannel:        &TickerFeedItem{},
	TickerBatchChannel:   &TickerBatchFeedItem{},
	Level2Channel:        &Level2FeedItem{},
	Channel("l2_data"):   &Level2FeedItem{},
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
		panic(msgType.Channel)
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

type CandlesFeedItem struct {
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

func (s *CandlesFeedItem) Seq() int64 {
	return s.SequenceNum
}

func (s *CandlesFeedItem) Clone() WebsocketMessage {
	return new(CandlesFeedItem)
}

type MarketTradesFeedItem struct {
	MessageDetails
	Events []*MarketTradesEvent `json:"events"`
}

func (s *MarketTradesFeedItem) Seq() int64 {
	return s.SequenceNum
}

func (s *MarketTradesFeedItem) Clone() WebsocketMessage {
	return new(MarketTradesFeedItem)
}

type MarketTradesEvent struct {
	Type   string        `json:"type"`
	Trades []*TradeEvent `json:"trades"`
}

type TradeEvent struct {
	TradeId   string    `json:"trade_id"`
	ProductId string    `json:"product_id"`
	Price     float64   `json:"price,string"`
	Size      float64   `json:"size,string"`
	Side      string    `json:"side"`
	Time      time.Time `json:"time"`
}

type StatusFeedItem struct {
	MessageDetails
	Events []*StatusEvent `json:"events"`
}

type StatusEvent struct {
	Type     string           `json:"type"`
	Products []*ProductUpdate `json:"products"`
}

type ProductUpdate struct {
	ProductType    ProductType `json:"product_type"`
	ProductId      string      `json:"id"`
	BaseCurrency   string      `json:"base_currency"`
	QuoteCurrency  string      `json:"quote_currency"`
	BaseIncrement  float64     `json:"base_increment,string"`
	QuoteIncrement float64     `json:"quote_increment,string"`
	DisplayName    string      `json:"display_name"`
	Status         string      `json:"status"`
	StatusMessage  string      `json:"status_message"`
	MinMarketFunds float64     `json:"min_market_funds,string"`
}

func (s *StatusFeedItem) Seq() int64 {
	return s.SequenceNum
}

func (s *StatusFeedItem) Clone() WebsocketMessage {
	return new(StatusFeedItem)
}

type TickerFeedItem struct {
	MessageDetails
	Events []*TickerEvent `json:"events"`
}

type TickerEvent struct {
	Type    string    `json:"type"`
	Tickers []*Ticker `json:"tickers"`
}

type Ticker struct {
	Type             string  `json:"type"`
	ProductId        string  `json:"product_id"`
	Price            float64 `json:"price,string"`
	Volume24H        float64 `json:"volume_24_h,string"`
	Low24H           float64 `json:"low_24_h,string"`
	High24H          float64 `json:"high_24_h,string"`
	Low52W           float64 `json:"low_52_w,string"`
	High52W          float64 `json:"high_52_w,string"`
	PercentChange24h float64 `json:"price_percent_chg_24_h,string"`
}

func (s *TickerFeedItem) Seq() int64 {
	return s.SequenceNum
}

func (s *TickerFeedItem) Clone() WebsocketMessage {
	return new(TickerFeedItem)
}

type TickerBatchFeedItem struct {
	MessageDetails
	Events []*TickerEvent `json:"events"`
}

func (s *TickerBatchFeedItem) Seq() int64 {
	return s.SequenceNum
}

func (s *TickerBatchFeedItem) Clone() WebsocketMessage {
	return new(TickerBatchFeedItem)
}

type Level2FeedItem struct {
	MessageDetails
	Events []*Level2Event `json:"events"`
}

type Level2Event struct {
	Type      string          `json:"type"`
	ProductId string          `json:"product_id"`
	Updates   []*Level2Update `json:"updates"`
}

type Level2Update struct {
	Side string `json:"side"`
	// This field is not populated by the coibnase API for some reason...
	EventTime   time.Time `json:"event_time"`
	PriceLevel  float64   `json:"price_level,string"`
	NewQuantity float64   `json:"new_quantity,string"`
}

func (s *Level2FeedItem) Seq() int64 {
	return s.SequenceNum
}

func (s *Level2FeedItem) Clone() WebsocketMessage {
	return new(Level2FeedItem)
}
