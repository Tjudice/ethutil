package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type AccountAuth struct {
	API_KEY        string `json:"api_key" yaml:"api_key"`
	API_PASSPHRASE string `json:"api_passphrase" yaml:"api_passphrase"`
	API_SECRET     string `json:"api_secret" yaml:"api_secret"`
}

func LoadAccount(filepath string) (*AccountAuth, error) {
	bts, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	out := &AccountAuth{}
	err = yaml.Unmarshal(bts, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func SignRequest(a *AccountAuth, endpoint string, r *http.Request) error {
	access_timestamp := time.Now().Unix()
	body, err := encodeBody(r)
	if err != nil {
		return err
	}
	message := strconv.FormatInt(access_timestamp, 10) + r.Method + endpoint + body
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
	Passphrase string `json:"passphrase"`
	Sig        string `json:"signature"`
}

func SignWebsocket(a *AccountAuth) (*SignedMessage, error) {
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
		Sig:        string(sig64),
	}
	return r, nil
}

func encodeBody(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", nil
	}
	body, err := r.GetBody()
	if err != nil {
		return "", err
	}
	bodyBts, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	if len(bodyBts) == 0 {
		return "", nil
	}
	marshalled, err := json.Marshal(bodyBts)
	if err != nil {
		return "", err
	}
	return string(marshalled), nil
}

// func Auth2(a *AccountAuth, endpoint string, r *http.Request) error {
// 	access_timestamp := time.Now().Unix()
// 	body, err := r.GetBody()
// 	if err != nil {
// 		return err
// 	}
// 	bodyBts, err := io.ReadAll(body)
// 	if err != nil {
// 		return err
// 	}
// 	marshalled, err := json.Marshal(bodyBts)
// 	if err != nil {
// 		return err
// 	}

// 	message := strconv.FormatInt(access_timestamp, 10) + endpoint + string(marshalled)

// 	r.Header.Set("ACCESS_KEY", a.API_KEY)

// 	h := hmac.New(sha256.New, []byte(a.API_SECRET))
// 	h.Write([]byte(message))

// 	signature := hex.EncodeToString(h.Sum(nil))

// 	r.Header.Set("ACCESS_SIGNATURE", signature)
// 	r.Header.Set("ACCESS_", nonce)

// 	return nil
// }