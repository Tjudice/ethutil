package coinbase_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/tjudice/ethutil/exchange/coinbase"
)

func getClient() *coinbase.Client {
	acctEnv := os.Getenv("AUTH_FILE_PATH")
	acc, _ := coinbase.LoadAccount(coinbase.ExchangeAuth, acctEnv)
	return coinbase.NewClient(acc)
}

func TestGetMarkets(t *testing.T) {
	cl := getClient()
	res, err := cl.GetMarkets(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res {
		log.Printf("%+v", r)
	}
}

func TestGetbook12(t *testing.T) {
	cl := getClient()
	book, err := cl.GetMarketBookLevel2(context.TODO(), "btc-usd")
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range book.Bids {
		log.Printf("%+v\n", b)
	}
}

func TestGetbook3(t *testing.T) {
	cl := getClient()
	book, err := cl.GetMarketBookLevel3(context.TODO(), "btc-usd")
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range book.Bids {
		log.Printf("%+v\n", b)
	}
}

func TestCandles(t *testing.T) {
	cl := getClient()
	cns, err := cl.GetMarketCandles(context.TODO(), "btc-usd", 60, int(time.Now().Unix()-600), int(time.Now().Unix()))
	if err != nil {
		t.Fatal(err)
	}
	for _, cn := range *cns {
		log.Printf("%+v", cn)
	}
}

func TestGetMarketStats(t *testing.T) {
	cl := getClient()
	stats, err := cl.GetMarketStats(context.TODO(), "btc-usd")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", stats)
}

func TestGetMarketTicker(t *testing.T) {
	cl := getClient()
	ticker, err := cl.GetMarketTicker(context.TODO(), "btc-usd")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(ticker)
}

func TestGetMarketTrades(t *testing.T) {
	cl := getClient()
	trades, err := cl.GetMarketTrades(context.TODO(), "btc-usd", 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, trade := range trades {
		log.Println(trade)
	}
}
