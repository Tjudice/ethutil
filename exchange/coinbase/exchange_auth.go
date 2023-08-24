package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"
)

type ExchangeAccountAuth struct {
	API_KEY        string `json:"api_key" yaml:"api_key"`
	API_PASSPHRASE string `json:"api_passphrase" yaml:"api_passphrase"`
	API_SECRET     string `json:"api_secret" yaml:"api_secret"`
}

func (a *ExchangeAccountAuth) SignRequest(requestPath string, r *http.Request) error {
	access_timestamp := time.Now().Unix()
	body, err := encodeBody(r)
	if err != nil {
		return err
	}
	message := strconv.FormatInt(access_timestamp, 10) + r.Method + requestPath + body
	decoded, err := base64.StdEncoding.DecodeString(a.API_SECRET)
	if err != nil {
		return err
	}
	secretHmac := hmac.New(sha256.New, decoded)
	_, err = secretHmac.Write([]byte(message))
	if err != nil {
		return err
	}
	sig := secretHmac.Sum(make([]byte, 0, secretHmac.Size()))
	sig64 := base64.StdEncoding.EncodeToString(sig)
	r.Header.Add("CB-ACCESS-KEY", a.API_KEY)
	r.Header.Add("CB-ACCESS-TIMESTAMP", strconv.FormatInt(access_timestamp, 10))
	r.Header.Add("CB-ACCESS-PASSPHRASE", a.API_PASSPHRASE)
	r.Header.Add("CB-ACCESS-SIGN", string(sig64))
	return nil
}

type SignedMessage struct {
	Key        string `json:"key"`
	Timestamp  string `json:"timestamp"`
	Passphrase string `json:"passphrase,omitempty"`
	Sig        string `json:"signature"`
}

func (a *ExchangeAccountAuth) SignWebsocketRequest(channels []string, products []string) (*SignedMessage, error) {
	access_timestamp := time.Now().Unix()
	message := strconv.FormatInt(access_timestamp, 10) + http.MethodGet + "/users/self/verify"
	decoded, err := base64.StdEncoding.DecodeString(a.API_SECRET)
	if err != nil {
		return nil, err
	}
	secretHmac := hmac.New(sha256.New, decoded)
	_, err = secretHmac.Write([]byte(message))
	if err != nil {
		return nil, err
	}
	sig := secretHmac.Sum(make([]byte, 0, secretHmac.Size()))
	sig64 := base64.StdEncoding.EncodeToString(sig)
	r := &SignedMessage{
		Key:        a.API_KEY,
		Timestamp:  strconv.FormatInt(access_timestamp, 10),
		Passphrase: a.API_PASSPHRASE,
		Sig:        sig64,
	}
	return r, nil
}
