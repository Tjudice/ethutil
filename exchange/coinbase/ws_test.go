package coinbase_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/tjudice/ethutil/exchange/coinbase"
)

func getWsClient() *coinbase.Client {
	acctEnv := os.Getenv("AUTH_FILE_PATH")
	acc, _ := coinbase.LoadAccount(acctEnv)
	return coinbase.NewClient(acc)
}

func TestSubscribeFull(t *testing.T) {
	cl := getWsClient()
	conn, err := cl.Subscribe(context.TODO(), []string{"MXC-USD"}, []any{"full"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Println(x)
		}
	}()
	time.Sleep(time.Minute)
}
