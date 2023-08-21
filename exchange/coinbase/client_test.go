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

func TestGetbook12(t *testing.T) {
	cl := coinbase.NewClient()
	book, err := cl.GetMarketBookLevel2(context.TODO(), "btc-usd")
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range book.Bids {
		log.Printf("%+v\n", b)
	}
}

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
