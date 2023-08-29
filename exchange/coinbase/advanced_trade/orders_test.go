package advanced_trade_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/tjudice/ethutil/exchange/coinbase"
	"github.com/tjudice/ethutil/exchange/coinbase/advanced_trade"
)

func getAdvancedTradeClient3() *advanced_trade.Client {
	acctEnv := os.Getenv("ADVANCED_TRADE_AUTH_FILE_PATH")
	acc, _ := coinbase.LoadAccount(coinbase.AdvancedTrade, acctEnv)
	return advanced_trade.NewClient(acc)
}

func TestGetOrders(t *testing.T) {
	cl := getAdvancedTradeClient3()
	res, err := cl.GetOrders(context.TODO(), &advanced_trade.OrderParams{
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
	res, err := cl.GetOrders(context.TODO(), &advanced_trade.OrderParams{
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
	res, err := cl.GetFills(context.TODO(), &advanced_trade.FillParams{
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
	res, err := cl.GetTransactionsSummary(context.TODO(), &advanced_trade.TransactionsSummaryParams{})
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", res)
}
