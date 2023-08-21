package coinbase_test

import (
	"context"
	"log"
	"testing"

	"github.com/tjudice/ethutil/exchange/coinbase"
)

// func TestGetMarkets(t *testing.T) {
// 	cl := coinbase.NewClient()
// 	res, err := cl.GetMarkets(context.TODO())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for _, r := range res {
// 		log.Printf("%+v", r)
// 	}
// }

// func TestGetbook12(t *testing.T) {
// 	cl := coinbase.NewClient()
// 	book, err := cl.GetMarketBookLevel2(context.TODO(), "btc-usd")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for _, b := range book.Bids {
// 		log.Printf("%+v\n", b)
// 	}
// }

// func TestGetbook3(t *testing.T) {
// 	cl := coinbase.NewClient()
// 	book, err := cl.GetMarketBookLevel3(context.TODO(), "btc-usd")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for _, b := range book.Bids {
// 		log.Printf("%+v\n", b)
// 	}
// }

// func TestCandles(t *testing.T) {
// 	cl := coinbase.NewClient()
// 	cns, err := cl.GetMarketCandles(context.TODO(), "btc-usd", 60, int(time.Now().Unix()-600), int(time.Now().Unix()))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for _, cn := range *cns {
// 		log.Printf("%+v", cn)
// 	}
// }

// func TestGetMarketStats(t *testing.T) {
// 	cl := coinbase.NewClient()
// 	stats, err := cl.GetMarketStats(context.TODO(), "btc-usd")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	log.Printf("%+v", stats)
// }

func TestGetMarketTicker(t *testing.T) {
	cl := coinbase.NewClient()
	ticker, err := cl.GetMarketTicker(context.TODO(), "btc-usd")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(ticker)
}

func TestGetMarketTrades(t *testing.T) {
	cl := coinbase.NewClient()
	trades, err := cl.GetMarketTrades(context.TODO(), "btc-usd", 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, trade := range trades {
		log.Println(trade)
	}
}
