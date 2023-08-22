package coinbase_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/tjudice/ethutil/exchange/coinbase"
)

func getAdvancedTradeClient() *coinbase.AdvancedTradeClient {
	acctEnv := os.Getenv("ADVANCED_TRADE_AUTH_FILE_PATH")
	acc, _ := coinbase.LoadAccount(coinbase.AdvancedTradeAuth, acctEnv)
	return coinbase.NewAdvancedTradeClient(acc)
}

func TestAdvancedTradeGetAccounts(t *testing.T) {
	cl := getAdvancedTradeClient()
	accts, err := cl.GetAccounts(context.TODO(), 10, "")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", accts)
}

func TestAdvancedTradeGetAccount(t *testing.T) {
	cl := getAdvancedTradeClient()
	accts, err := cl.GetAccounts(context.TODO(), 10, "")
	if err != nil {
		t.Fatal(err)
	}
	uuid := accts.Accounts[0].UUID
	acct, err := cl.GetAccount(context.TODO(), uuid)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", acct)
}
