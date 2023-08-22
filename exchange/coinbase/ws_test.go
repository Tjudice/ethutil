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
	time.Sleep(10 * time.Second)
}

func TestSubscribeTicker(t *testing.T) {
	cl := getWsClient()
	conn, err := cl.Subscribe(context.TODO(), []string{"MXC-USD"}, []any{"ticker"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Println(x)
		}
	}()
	time.Sleep(10 * time.Second)
}

func TestSubscribeStatus(t *testing.T) {
	cl := getWsClient()
	conn, err := cl.Subscribe(context.TODO(), []string{"MXC-USD"}, []any{"status"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			x := <-conn.C()
			log.Printf("%+v", x)
		}
	}()
	time.Sleep(10 * time.Second)
}

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
	time.Sleep(10 * time.Second)
	log.Println(processed)
}

func TestSubscribeTickerBatch(t *testing.T) {
	cl := getWsClient()
	conn, err := cl.Subscribe(context.TODO(), []string{"BTC-USD"}, []any{"ticker_batch"})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			log.Println(<-conn.C())
		}
	}()
	time.Sleep(10 * time.Second)
}

func TestSubscribeRFQMatches(t *testing.T) {
	cl := getWsClient()
	conn, err := cl.Subscribe(context.TODO(), nil, []any{
		struct {
			Name string `json:"name"`
		}{
			Name: "rfq_matches"},
	})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			log.Println(<-conn.C())
		}
	}()
	time.Sleep(10 * time.Second)
}

func TestSusbcribeLevel2(t *testing.T) {
	cl := getWsClient()
	conn, err := cl.Subscribe(context.TODO(), []string{"BTC-USD"}, []any{"level2"})
	if err != nil {
		t.Fatal(err)
	}
	cnt := 0
	go func() {
		for {
			<-conn.C()
			cnt = cnt + 1
		}
	}()
	time.Sleep(10 * time.Second)
	log.Println(cnt)
}
