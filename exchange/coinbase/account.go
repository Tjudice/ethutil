package coinbase

import (
	"os"

	"github.com/tjudice/ethutil/exchange/coinbase/advanced_trade"
	"github.com/tjudice/ethutil/exchange/coinbase/auth"
	"github.com/tjudice/ethutil/exchange/coinbase/exchange"

	"gopkg.in/yaml.v2"
)

type AuthenticationType int

var (
	Exchange      AuthenticationType = 1
	AdvancedTrade AuthenticationType = 2
	OAUTH2        AuthenticationType = 3
)

func LoadAccount(accountType AuthenticationType, filepath string) (auth.Authenticator, error) {
	bts, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var out auth.Authenticator
	switch accountType {
	case Exchange:
		out = &exchange.Auth{}
	case AdvancedTrade:
		out = &advanced_trade.Auth{}
	case OAUTH2:
		panic("not implemented")
	}
	err = yaml.Unmarshal(bts, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
