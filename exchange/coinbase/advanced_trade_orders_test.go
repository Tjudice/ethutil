package coinbase_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/tjudice/ethutil/exchange/coinbase"
)

func getAdvancedTradeClient3() *coinbase.AdvancedTradeClient {
	acctEnv := os.Getenv("ADVANCED_TRADE_AUTH_FILE_PATH")
	acc, _ := coinbase.LoadAccount(coinbase.AdvancedTradeAuth, acctEnv)
	return coinbase.NewAdvancedTradeClient(acc)
}

func TestGetOrders(t *testing.T) {
	cl := getAdvancedTradeClient3()
	res, err := cl.GetOrders(context.TODO(), &coinbase.OrderParams{
		Limit: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(res))
}
