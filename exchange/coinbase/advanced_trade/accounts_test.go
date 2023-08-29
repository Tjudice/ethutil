package advanced_trade_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/tjudice/ethutil/exchange/coinbase"
	"github.com/tjudice/ethutil/exchange/coinbase/advanced_trade"
)

func getAdvancedTradeClient() *advanced_trade.Client {
	acctEnv := os.Getenv("ADVANCED_TRADE_AUTH_FILE_PATH")
	acc, _ := coinbase.LoadAccount(coinbase.AdvancedTrade, acctEnv)
	return advanced_trade.NewClient(acc)
}

func TestGetAccounts(t *testing.T) {
	cl := getAdvancedTradeClient()
	accts, err := cl.GetAccounts(context.TODO(), 10, "")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", accts)
	log.Println(len(accts.Accounts))
}

func TestGetAccountsCursor(t *testing.T) {
	cl := getAdvancedTradeClient()
	accts, err := cl.GetAccounts(context.TODO(), 10, "")
	if err != nil {
		t.Fatal(err)
	}
	accts2, err := cl.GetAccounts(context.TODO(), 10, accts.Cursor)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(accts.Accounts[0], accts2.Accounts[0])
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
