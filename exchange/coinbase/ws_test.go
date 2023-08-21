package coinbase_test

import (
	"context"
	"testing"
)

func TestSubscribeFull(t *testing.T) {
	cl := getClient()
	_, err := cl.Subscribe(context.TODO(), []string{"MXC-USD"}, []any{"full"})
	if err != nil {
		t.Fatal(err)
	}
}
