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
	acc, err := coinbase.LoadAccount(acctEnv)
	if err != nil {
		panic(err)
	}
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

// func TestSubscribeStatus(t *testing.T) {
// 	cl := getWsClient()
// 	conn, err := cl.Subscribe(context.TODO(), []string{"MXC-USD"}, []any{"status"})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	go func() {
// 		for {
// 			x := <-conn.C()
// 			log.Printf("%+v", x)
// 		}
// 	}()
// 	time.Sleep(time.Minute)
// }

func TestSubscribeLevel3(t *testing.T) {
	cl := getWsClient()
	conn, err := cl.Subscribe(context.TODO(), []string{"BTC-USD"}, []any{"level3"})
	if err != nil {
		t.Fatal(err)
	}
	processed := 0
	go func() {
		for {
			<-conn.C()
			processed += 1
		}
	}()
	time.Sleep(time.Minute)
	log.Println(processed)
}
