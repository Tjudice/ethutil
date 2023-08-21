package coinbase_test

import (
	"context"
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

// func TestSubscribeFull(t *testing.T) {
// 	cl := getWsClient()
// 	conn, err := cl.Subscribe(context.TODO(), []string{"MXC-USD"}, []any{"full"})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	go func() {
// 		for {
// 			x := <-conn.C()
// 			log.Println(x)
// 		}
// 	}()
// 	time.Sleep(time.Minute)
// }

// func TestSubscribeTicker(t *testing.T) {
// 	cl := getWsClient()
// 	conn, err := cl.Subscribe(context.TODO(), []string{"MXC-USD"}, []any{"ticker"})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	go func() {
// 		for {
// 			x := <-conn.C()
// 			log.Println(x)
// 		}
// 	}()
// 	time.Sleep(time.Minute)
// }

func TestSubscribeStatus(t *testing.T) {
	cl := getWsClient()
	conn, err := cl.Subscribe(context.TODO(), []string{"MXC-USD"}, []any{"status"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			<-conn.C()
		}
	}()
	time.Sleep(time.Minute)
}
