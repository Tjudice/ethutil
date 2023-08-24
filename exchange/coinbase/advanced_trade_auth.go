package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gfx.cafe/open/ghost/hexutil"
)

type AdvancedTradeAuthenticator struct {
	API_KEY    string `json:"api_key" yaml:"api_key"`
	API_SECRET string `json:"api_secret" yaml:"api_secret"`
}

func (a *AdvancedTradeAuthenticator) SignRequest(requestPath string, r *http.Request) error {
	access_timestamp := time.Now().Unix()
	body, err := encodeBody(r)
	if err != nil {
		return err
	}
	message := strconv.FormatInt(access_timestamp, 10) + r.Method + requestPath + body
	secretHmac := hmac.New(sha256.New, []byte(a.API_SECRET))
	_, err = secretHmac.Write([]byte(message))
	if err != nil {
		return err
	}
	sig := secretHmac.Sum(make([]byte, 0, secretHmac.Size()))
	r.Header.Add("CB-ACCESS-KEY", a.API_KEY)
	r.Header.Add("CB-ACCESS-TIMESTAMP", strconv.FormatInt(access_timestamp, 10))
	r.Header.Add("CB-ACCESS-SIGN", hexutil.Bytes(sig).String()[2:])
	return nil
}

func (a *AdvancedTradeAuthenticator) SignWebsocketRequest(channels []string, products []string) (*SignedMessage, error) {
	access_timestamp := time.Now().Unix()
	channelStr := strings.Join(channels, "") + strings.Join(products, ",")
	message := strconv.FormatInt(access_timestamp, 10) + channelStr
	secretHmac := hmac.New(sha256.New, []byte(a.API_SECRET))
	_, err := secretHmac.Write([]byte(message))
	if err != nil {
		return nil, err
	}
	sig := secretHmac.Sum(make([]byte, 0, secretHmac.Size()))
	r := &SignedMessage{
		Key:       a.API_KEY,
		Timestamp: strconv.FormatInt(access_timestamp, 10),
		Sig:       hexutil.Bytes(sig).String()[2:],
	}
	return r, nil
}
