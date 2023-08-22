package coinbase

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson/fastfloat"
)

type Conn struct {
	ch   chan WebsocketMessage
	conn *websocket.Conn
	auth *AccountAuth
}

const (
	COINBASE_EXCHANGE_WSS_URL_FEED   = "wss://ws-feed.exchange.coinbase.com"
	COINBASE_EXCHANGE_WSS_URL_DIRECT = "wss://ws-direct.exchange.coinbase.com"
)

type SubscribeMsg struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []any    `json:"channels"`
}

type SignedSubscribeMsg struct {
	*SubscribeMsg
	*SignedMessage
}

type WebsocketMessage interface {
	Seq() int64
	Clone() WebsocketMessage
}

func (c *Client) Subscribe(ctx context.Context, products []string, channels []any) (*Conn, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, COINBASE_EXCHANGE_WSS_URL_FEED, nil)
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

func (c *Conn) Unsubscribe(ctx context.Context, products []string, channels []any) error {
	msg := &SubscribeMsg{
		Type:       "unsubscribe",
		ProductIds: products,
		Channels:   channels,
	}
	return c.conn.WriteJSON(msg)
}

func (c *Conn) Listen() error {
	defer c.conn.Close()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		parsed, err := parseMessage(msg)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		c.ch <- parsed
	}
}

func (c *Conn) C() <-chan WebsocketMessage {
	return c.ch
}

type MessageType struct {
	Type string `json:"type"`
}

var messageTypeChoice = map[string]WebsocketMessage{
	"subscriptions": &Subscriptions{},
	"done":          &Done{},
	"received":      &Received{},
	"open":          &Open{},
	"match":         &Match{},
	"ticker":        &WsTicker{},
	"status":        &Status{},
	"change":        &Change{},
	"activate":      &Activate{},
	"level3":        &Level3{},
}

func parseMessage(bts []byte) (WebsocketMessage, error) {
	if len(bts) == 0 {
		return nil, fmt.Errorf("zero length message")
	}
	if bts[0] == 123 {
		return parseJson(bts)
	}
	return parseArray(bts)
}

func parseJson(bts []byte) (WebsocketMessage, error) {
	msgType := MessageType{}
	err := json.Unmarshal(bts, &msgType)
	if err != nil {
		return nil, err
	}
	schem, ok := messageTypeChoice[msgType.Type]
	if !ok {
		log.Println(msgType.Type)
		log.Println(string(bts))
		panic("a")
	}
	typedMessage := schem.Clone()
	err = json.Unmarshal(bts, &typedMessage)
	if err != nil {
		log.Println(msgType.Type)
		log.Println(string(bts))
		return nil, err
	}
	return typedMessage, nil
}

type ArrayMessage interface {
	Populate(decodedArray []json.RawMessage) error
	WebsocketMessage
}

var messageTypeArray = map[string]ArrayMessage{
	"open":   &OpenL3{},
	"change": &ChangeL3{},
	"noop":   &NoopL3{},
	"match":  &MatchL3{},
	"done":   &DoneL3{},
}

func parseArray(bts []byte) (WebsocketMessage, error) {
	var arr []json.RawMessage
	err := json.Unmarshal(bts, &arr)
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("zero length message array")
	}
	var messageType string
	if err := json.Unmarshal(arr[0], &messageType); err != nil {
		return nil, err
	}
	typed, ok := messageTypeArray[messageType]
	if !ok {
		return nil, fmt.Errorf("unknown message type: %s", messageType)
	}
	parseInto := typed.Clone()
	populator := parseInto.(ArrayMessage)
	err = populator.Populate(arr[1:])
	if err != nil {
		return nil, err
	}
	return parseInto, nil
}

type Subscriptions struct {
	Channels []*ChannelSubscription `json:"channels"`
}

type ChannelSubscription struct {
	Name       string   `json:"name"`
	ProductIds []string `json:"product_ids"`
}

func (s *Subscriptions) Seq() int64 {
	return 0
}

func (s *Subscriptions) Clone() WebsocketMessage {
	return new(Subscriptions)
}

type Done struct {
	OrderId       string    `json:"order_id"`
	Reason        string    `json:"reason"`
	Price         float64   `json:"price,string"`
	RemainingSize float64   `json:"remaining_size,string"`
	Side          string    `json:"side"`
	ProductId     string    `json:"product_id"`
	Time          time.Time `json:"time"`
	Sequence      int64     `json:"sequence"`
}

func (s *Done) Seq() int64 {
	return s.Sequence
}

func (s *Done) Clone() WebsocketMessage {
	return new(Done)
}

type Received struct {
	OrderId   string    `json:"order_id"`
	OrderType string    `json:"order_type"`
	Size      float64   `json:"size,string"`
	Price     float64   `json:"price,string"`
	Side      string    `json:"side"`
	ProductId string    `json:"product_id"`
	Time      time.Time `json:"time"`
	Sequence  int64     `json:"sequence"`
}

func (s *Received) Seq() int64 {
	return s.Sequence
}

func (s *Received) Clone() WebsocketMessage {
	return new(Received)
}

type Open struct {
	OrderId       string    `json:"order_id"`
	RemainingSize float64   `json:"remaining_size,string"`
	Price         float64   `json:"price,string"`
	Side          string    `json:"side"`
	ProductId     string    `json:"product_id"`
	Time          time.Time `json:"time"`
	Sequence      int64     `json:"sequence"`
}

func (s *Open) Seq() int64 {
	return s.Sequence
}

func (s *Open) Clone() WebsocketMessage {
	return new(Open)
}

type Match struct {
	TradeId      int64     `json:"trade_id"`
	MakerOrderId string    `json:"maker_order_id"`
	TakerOrderId string    `json:"taker_order_id"`
	Size         float64   `json:"size,string"`
	Price        float64   `json:"price,string"`
	Side         string    `json:"side"`
	ProductId    string    `json:"product_id"`
	Time         time.Time `json:"time"`
	Sequence     int64     `json:"sequence"`
}

func (s *Match) Seq() int64 {
	return s.Sequence
}

func (s *Match) Clone() WebsocketMessage {
	return new(Match)
}

type WsTicker struct {
	ProductId   string    `json:"product_id"`
	Price       float64   `json:"price,string"`
	Open24H     float64   `json:"open_24h,string"`
	Volume24H   float64   `json:"volume_24h,string"`
	Low24H      float64   `json:"low_24h,string"`
	High24H     float64   `json:"high_24h,string"`
	Volume30D   float64   `json:"volume_30d,string"`
	BestBid     float64   `json:"best_bid,string"`
	BestBidSize float64   `json:"best_bid_size,string"`
	BestAsk     float64   `json:"best_ask,string"`
	BestAskSize float64   `json:"best_ask_size,string"`
	Side        string    `json:"side"`
	Time        time.Time `json:"time"`
	TradeId     int64     `json:"trade_id"`
	LastSize    float64   `json:"last_size,string"`
	Sequence    int64     `json:"sequence"`
}

func (s *WsTicker) Seq() int64 {
	return s.Sequence
}

func (s *WsTicker) Clone() WebsocketMessage {
	return new(WsTicker)
}

type Currency struct {
	Id                string     `json:"id"`
	Name              string     `json:"name"`
	MinSize           float64    `json:"min_size,string"`
	Status            string     `json:"status"`
	FundingAccountId  string     `json:"funding_account_id"`
	StatusMessage     string     `json:"status_message"`
	MaxPrecision      float64    `json:"max_precision,string"`
	ConvertibleTo     []any      `json:"convertible_to"`
	Details           *Details   `json:"details"`
	DefaultNetwork    string     `json:"default_network"`
	SupportedNetworks []*Network `json:"supported_networks"`
}

type Details struct {
	Type                  string   `json:"type"`
	Symbol                string   `json:"symbol"`
	Networkconfirmations  int      `json:"network_confirmations"`
	SortOrder             int      `json:"sort_order"`
	CryptoAddressLink     string   `json:"crypto_address_link"`
	CryptoTransactionLink string   `json:"crypto_transaction_link"`
	PushPaymentMethods    []string `json:"push_payment_methods"`
	MinWithdrawalAmount   float64  `json:"min_withdrawal_amount"`
	MaxWithdrawalAmount   float64  `json:"max_withdrawal_amount"`
}

type Network struct {
	Id                    string  `json:"id"`
	Name                  string  `json:"name"`
	Status                string  `json:"status"`
	ContractAddress       string  `json:"contract_address"`
	CryptoAddressLink     string  `json:"crypto_address_link"`
	CryptoTransactionLink string  `json:"crypto_transaction_link"`
	MinWithdrawalAmount   float64 `json:"min_withdrawal_amount"`
	MaxWithdrawalAmount   float64 `json:"max_withdrawal_amount"`
	Networkconfirmations  int     `json:"network_confirmations"`
	ProcessingTimeSeconds int     `json:"processing_time_seconds"`
	NetworkMap            any     `json:"network_map"`
}

type Status struct {
	Currencies []*Currency `json:"currencies"`
	Markets    []*Market   `json:"products"`
}

func (s *Status) Seq() int64 {
	return 0
}

func (s *Status) Clone() WebsocketMessage {
	return new(Status)
}

type Change struct {
	Reason    string    `json:"reason"`
	Time      time.Time `json:"time"`
	OrderId   string    `json:"order_id"`
	Side      string    `json:"side"`
	ProductId string    `json:"product_id"`
	OldSize   float64   `json:"old_size,string"`
	NewSize   float64   `json:"new_size,string"`
	OldPrice  float64   `json:"old_price,string"`
	NewPrice  float64   `json:"new_price,string"`
	Sequence  int64     `json:"sequence"`
}

func (s *Change) Seq() int64 {
	return s.Sequence
}

func (s *Change) Clone() WebsocketMessage {
	return new(Change)
}

type Activate struct {
	ProductId string  `json:"product_id"`
	Timestamp float64 `json:"timestamp,string"`
	UserId    int     `json:"user_id,string"`
	ProfileId string  `json:"profile_id"`
	OrderId   string  `json:"order_id"`
	StopType  string  `json:"stop_type"`
	Side      string  `json:"side"`
	StopPrice float64 `json:"stop_price,string"`
	Size      float64 `json:"size,string"`
	Funds     float64 `json:"funds,string"`
	Private   bool    `json:"private"`
}

func (s *Activate) Seq() int64 {
	return 0
}

func (s *Activate) Clone() WebsocketMessage {
	return new(Activate)
}

type Snapshot struct{}

type L2Update struct{}

type Level3Schema struct {
	Change []string `json:"change"`
	Done   []string `json:"done"`
	Match  []string `json:"match"`
	Noop   []string `json:"noop"`
	Open   []string `json:"open"`
}

type Level3 struct {
	Schema *Level3Schema
}

func (s *Level3) Seq() int64 {
	return 0
}

func (s *Level3) Clone() WebsocketMessage {
	return new(Level3)
}

type OpenL3 struct {
	ProductId string    `json:"product_id"`
	Sequence  int64     `json:"sequence"`
	OrderId   string    `json:"order_id"`
	Side      string    `json:"side"`
	Price     float64   `json:"price,string"`
	Size      float64   `json:"size,string"`
	Time      time.Time `json:"time"`
}

func (s *OpenL3) Seq() int64 {
	return s.Sequence
}

func (s *OpenL3) Clone() WebsocketMessage {
	return new(OpenL3)
}

func (s *OpenL3) Populate(r []json.RawMessage) (err error) {
	if len(r) != 7 {
		return fmt.Errorf("open level 3: incorrect array length")
	}
	s.ProductId = string(r[0])
	s.Sequence = int64(binary.BigEndian.Uint64(r[1]))
	s.OrderId = string(r[2])
	s.Side = string(r[3])
	s.Price, err = fastfloat.Parse(string(r[4][1 : len(r[4])-1]))
	if err != nil {
		return err
	}
	s.Size, err = fastfloat.Parse(string(r[5][1 : len(r[5])-1]))
	if err != nil {
		return err
	}
	return s.Time.UnmarshalJSON(r[6])
}

type MatchL3 struct {
	ProductId    string    `json:"product_id"`
	Sequence     int64     `json:"sequence"`
	MakerOrderId string    `json:"maker_order_id"`
	TakerOrderId string    `json:"taker_order_id"`
	Price        float64   `json:"price,string"`
	Size         float64   `json:"size,string"`
	Time         time.Time `json:"time"`
}

func (s *MatchL3) Seq() int64 {
	return s.Sequence
}

func (s *MatchL3) Clone() WebsocketMessage {
	return new(MatchL3)
}

func (s *MatchL3) Populate(r []json.RawMessage) (err error) {
	if len(r) != 7 {
		return fmt.Errorf("match level 3: incorrect array length")
	}
	s.ProductId = string(r[0])
	s.Sequence = int64(binary.BigEndian.Uint64(r[1]))
	s.MakerOrderId = string(r[2])
	s.TakerOrderId = string(r[3])
	s.Price, err = fastfloat.Parse(string(r[4][1 : len(r[4])-1]))
	if err != nil {
		return err
	}
	s.Size, err = fastfloat.Parse(string(r[5][1 : len(r[5])-1]))
	if err != nil {
		return err
	}
	return s.Time.UnmarshalJSON(r[6])
}

type ChangeL3 struct {
	ProductId string    `json:"product_id"`
	Sequence  int64     `json:"sequence"`
	OrderId   string    `json:"order_id"`
	Price     float64   `json:"price,string"`
	Size      float64   `json:"size,string"`
	Time      time.Time `json:"time"`
}

func (s *ChangeL3) Seq() int64 {
	return s.Sequence
}

func (s *ChangeL3) Clone() WebsocketMessage {
	return new(ChangeL3)
}

func (s *ChangeL3) Populate(r []json.RawMessage) (err error) {
	if len(r) != 6 {
		return fmt.Errorf("change level 3: incorrect array length")
	}
	s.ProductId = string(r[0])
	s.Sequence = int64(binary.BigEndian.Uint64(r[1]))
	s.OrderId = string(r[2])
	s.Price, err = fastfloat.Parse(string(r[3][1 : len(r[3])-1]))
	if err != nil {
		return err
	}
	s.Size, err = fastfloat.Parse(string(r[4][1 : len(r[4])-1]))
	if err != nil {
		return err
	}
	return s.Time.UnmarshalJSON(r[5])
}

type NoopL3 struct {
	ProductId string    `json:"product_id"`
	Sequence  int64     `json:"sequence"`
	Time      time.Time `json:"time"`
}

func (s *NoopL3) Seq() int64 {
	return s.Sequence
}

func (s *NoopL3) Clone() WebsocketMessage {
	return new(NoopL3)
}

func (s *NoopL3) Populate(r []json.RawMessage) error {
	if len(r) != 3 {
		return fmt.Errorf("noop level 3: incorrect array length")
	}
	s.ProductId = string(r[0])
	s.Sequence = int64(binary.BigEndian.Uint64(r[1]))
	return s.Time.UnmarshalJSON(r[2])
}

type DoneL3 struct {
	ProductId string    `json:"product_id"`
	Sequence  int64     `json:"sequence"`
	OrderId   string    `json:"order_id"`
	Time      time.Time `json:"time"`
}

func (s *DoneL3) Seq() int64 {
	return s.Sequence
}

func (s *DoneL3) Clone() WebsocketMessage {
	return new(DoneL3)
}

func (s *DoneL3) Populate(r []json.RawMessage) error {
	if len(r) != 4 {
		return fmt.Errorf("done level 3: incorrect array length")
	}
	s.ProductId = string(r[0])
	s.Sequence = int64(binary.BigEndian.Uint64(r[1]))
	s.OrderId = string(r[2])
	return s.Time.UnmarshalJSON(r[3])
}
