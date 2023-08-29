package advanced_trade_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/tjudice/ethutil/exchange/coinbase"
	"github.com/tjudice/ethutil/exchange/coinbase/advanced_trade"
)

func getAdvancedTradeClient4() *advanced_trade.Client {
	acctEnv := os.Getenv("ADVANCED_TRADE_AUTH_FILE_PATH")
	acc, _ := coinbase.LoadAccount(coinbase.AdvancedTrade, acctEnv)
	return advanced_trade.NewClient(acc)
}

func TestSubscribeHeartbeats(t *testing.T) {
	cl := getAdvancedTradeClient4()
	conn, err := cl.Subscribe(context.TODO(), 10, advanced_trade.HeartbeatsChannel, []string{"BTC-USD", "ETH-USD"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Printf("%+v", x)
		}
	}()
	time.Sleep(10 * time.Second)
}

func TestSubscribeCandles(t *testing.T) {
	cl := getAdvancedTradeClient4()
	conn, err := cl.Subscribe(context.TODO(), 10, advanced_trade.CandlesChannel, []string{"BTC-USD", "ETH-USD"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Printf("%+v", x)
		}
	}()
	time.Sleep(10 * time.Second)
}

func TestSubscribeMarketTrades(t *testing.T) {
	cl := getAdvancedTradeClient4()
	conn, err := cl.Subscribe(context.TODO(), 10, advanced_trade.MarketTradesChannel, []string{"BTC-USD", "ETH-USD"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Printf("%+v", x)

		}
	}()
	time.Sleep(10 * time.Second)
}

func TestSubscribeStatus(t *testing.T) {
	cl := getAdvancedTradeClient4()
	conn, err := cl.Subscribe(context.TODO(), 10, advanced_trade.StatusChannel, []string{"BTC-USD", "ETH-USD"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Printf("%T %+v", x, x)

		}
	}()
	time.Sleep(10 * time.Second)
}

func TestSubscribeTicker(t *testing.T) {
	cl := getAdvancedTradeClient4()
	conn, err := cl.Subscribe(context.TODO(), 10, advanced_trade.TickerChannel, []string{"BTC-USD", "ETH-USD"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Printf("%T %+v", x, x)

		}
	}()
	time.Sleep(10 * time.Second)
}

func TestSubscribeTickerBatch(t *testing.T) {
	cl := getAdvancedTradeClient4()
	conn, err := cl.Subscribe(context.TODO(), 10, advanced_trade.TickerBatchChannel, []string{"BTC-USD", "ETH-USD"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Printf("%T %+v", x, x)

		}
	}()
	time.Sleep(10 * time.Second)
}

func TestSubscribeLevel2(t *testing.T) {
	cl := getAdvancedTradeClient4()
	conn, err := cl.Subscribe(context.TODO(), 10, advanced_trade.Level2Channel, []string{"BTC-USD", "ETH-USD"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Printf("%T %+v", x, x)

		}
	}()
	time.Sleep(10 * time.Second)
}
