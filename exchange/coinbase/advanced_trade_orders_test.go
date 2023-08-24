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
		Limit: 3,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range res.Orders {
		log.Printf("%+v", o)
	}
}

func TestGetOrder(t *testing.T) {
	cl := getAdvancedTradeClient3()
	res, err := cl.GetOrders(context.TODO(), &coinbase.OrderParams{
		Limit: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Orders) == 0 {
		t.Fatalf("user has no orders")
	}
	singleOrder, err := cl.GetOrder(context.TODO(), res.Orders[0].OrderId)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", singleOrder)
}

func TestGetFills(t *testing.T) {
	cl := getAdvancedTradeClient3()
	res, err := cl.GetFills(context.TODO(), &coinbase.FillParams{
		Limit: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Fills) == 0 {
		t.Fatalf("user has no fills")
	}
	log.Println(res.Fills[0])
}

func TestGetTransactionsSummary(t *testing.T) {
	cl := getAdvancedTradeClient3()
	res, err := cl.GetTransactionsSummary(context.TODO(), &coinbase.TransactionsSummaryParams{})
	if err != nil {
		t.Fatal(err)
	}
	_ = res
}
