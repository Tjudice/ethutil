package coinbase_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/tjudice/ethutil/exchange/coinbase"
)

func TestSubscribe(t *testing.T) {
	acctEnv := os.Getenv("AUTH_FILE_PATH")
	acc, err := coinbase.LoadAccount(acctEnv)
	if err != nil {
		t.Fatal(err)
	}
	_, err = coinbase.Subscribe(context.TODO(), acc, []string{"MXC-USD"}, []any{"full"})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Minute)
}
