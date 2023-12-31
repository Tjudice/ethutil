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

func getAdvancedTradeClient2() *advanced_trade.Client {
	acctEnv := os.Getenv("ADVANCED_TRADE_AUTH_FILE_PATH")
	acc, _ := coinbase.LoadAccount(coinbase.AdvancedTrade, acctEnv)
	return advanced_trade.NewClient(acc)
}

func TestGetBestBidAsk(t *testing.T) {
	cl := getAdvancedTradeClient2()
	res, err := cl.GetBestBidAsk(context.TODO(), []string{"BTC-USD", "ETH-USD"})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(res.PriceBooks[0].Bids[0])
}

func TestGetOrderbook(t *testing.T) {
	cl := getAdvancedTradeClient2()
	res, err := cl.GetOrderbook(context.TODO(), "NEST-USD", 10)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(res)
}

func TestGetMarketsAdvanced(t *testing.T) {
	cl := getAdvancedTradeClient2()
	res, err := cl.GetMarkets(context.TODO(), &advanced_trade.GetMarketParams{
		Limit:  100,
		Offset: 50,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(res.Products[20])
}

func TestGetMarketAdvanced(t *testing.T) {
	cl := getAdvancedTradeClient2()
	res, err := cl.GetMarket(context.TODO(), "BTC-USD")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(res)
}

func TestGetMarketCandles(t *testing.T) {
	cl := getAdvancedTradeClient2()
	currUnix := time.Now().Unix()
	res, err := cl.GetCandles(context.TODO(), "BTC-USD", advanced_trade.CANDLE_GRANULARITY_1_MINUTE, currUnix-60*100, currUnix)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range res.Candles {
		log.Printf("%+v", c)
	}
}

func TestGetAdvancedMarketTrades(t *testing.T) {
	cl := getAdvancedTradeClient2()
	res, err := cl.GetMarketTrades(context.TODO(), "BTC-USD", 5)
	if err != nil {
		t.Fatal(err)
	}
	for _, trade := range res.Trades {
		log.Printf("%+v", trade)
	}
}
