package coinbase_test

import (
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/tjudice/ethutil/exchange/coinbase"
)

func TestGetAuthentication(t *testing.T) {
	acctEnv := os.Getenv("AUTH_FILE_PATH")
	acc, err := coinbase.LoadAccount(acctEnv)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(acc)
}

func TestAuthenticateRequest(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://api.exchange.coinbase.com/accounts", nil)
	acctEnv := os.Getenv("AUTH_FILE_PATH")
	acc, err := coinbase.LoadAccount(acctEnv)
	if err != nil {
		t.Fatal(err)
	}
	err = coinbase.SignRequest(acc, "/accounts", req)
	if err != nil {
		t.Fatal(err)
	}
	cl := http.DefaultClient
	res, err := cl.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(res.Body)
	log.Println(string(b))
}
