package coinbase_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/tjudice/ethutil/exchange/coinbase"
)

func getAdvancedTradeClient2() *coinbase.AdvancedTradeClient {
	acctEnv := os.Getenv("ADVANCED_TRADE_AUTH_FILE_PATH")
	acc, _ := coinbase.LoadAccount(coinbase.AdvancedTradeAuth, acctEnv)
	return coinbase.NewAdvancedTradeClient(acc)
}

// func TestGetBestBidAsk(t *testing.T) {
// 	cl := getAdvancedTradeClient2()
// 	res, err := cl.GetBestBidAsk(context.TODO(), []string{"BTC-USD", "ETH-USD"})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	log.Println(res.PriceBooks[0].Bids[0])
// }

// func TestGetOrderbook(t *testing.T) {
// 	cl := getAdvancedTradeClient2()
// 	res, err := cl.GetOrderbook(context.TODO(), "NEST-USD", 10)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	log.Println(res)
// }

func TestGetMarketsAdvanced(t *testing.T) {
	cl := getAdvancedTradeClient2()
	res, err := cl.GetMarkets(context.TODO(), &coinbase.GetMarketParams{
		Limit:  100,
		Offset: 50,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(res))
}
