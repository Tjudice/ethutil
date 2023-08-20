package coinbase_test

import (
	"context"
	"log"
	"testing"

	"github.com/tjudice/ethutil/exchange/coinbase"
)

func TestGetMarkets(t *testing.T) {
	cl := coinbase.NewClient()
	res, err := cl.GetMarkets(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range res {
		log.Printf("%+v", r)
	}
}
